package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"time"

	log "github.com/mr-dai/static-server/log"
	mime "github.com/mr-dai/static-server/mime"
)

var logger = log.NewLogger(log.INFO, "static-server", "")

var enableCache bool
var enableList bool
var port int
var dir string

func main() {
	// Initialize
	flag.BoolVar(&enableCache, "cache", false, "enable HTTP caching")
	flag.BoolVar(&enableList, "list", false, "enable listing on directory when index.html is missed")
	flag.IntVar(&port, "port", 8080, "specify the port number to listen on")
	flag.StringVar(&dir, "dir", "", "path of directory to server, default to be current working directory")
	flag.Parse()

	if dir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			logger.Fatal(err.Error())
		}
		dir = pwd
	}

	logger.Debugf("Received configuration { port = %d, dir = %s, enableList = %t, enableCache = %t }",
		port, dir, enableList, enableCache)
	logger.Infof("Server is listening on %d...", port)

	// Setup HTTP server
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(fmt.Sprintf(":%d", port), http.DefaultServeMux)

	// TODO Wait on interrupt signal

	logger.Info("Received interrupt signal, exiting")
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	wrapped := Wrap(&w)
	handleRequestWrapped(&wrapped, req)
	logger.Infof("%s %s %s %d %d", req.Method, req.RequestURI, req.Proto, wrapped.code, wrapped.length)
}

// handleRequest is the handler function for the HTTP server
func handleRequestWrapped(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed) // Method not allowed
		return
	}
	path := filepath.Join(dir, req.RequestURI)
	logger.Tracef("Checking %s", path)
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) { // Append `.html` suffix and try again
		path += ".html"
		fileInfo, err = os.Stat(path)
	}
	if err != nil {
		handleFileError(err, path, w)
		return
	}

	if fileInfo.IsDir() {
		logger.Tracef("%s is directory", path)
		handleRequestForDirectory(w, req, path, fileInfo)
	} else {
		logger.Tracef("%s is file", path)
		handleRequestForFile(w, req, path, fileInfo)
	}
}

// handleRequestForFile handles the request towards the local file with the given path
func handleRequestForFile(w http.ResponseWriter, req *http.Request, path string, fileInfo os.FileInfo) {
	file, err := os.Open(path)
	if err != nil {
		handleFileError(err, path, w)
		return
	}
	defer file.Close()

	w.Header().Add("Content-Type", mime.GetMimeByFilename(file.Name()))
	if enableCache {
		w.Header().Add("Last-Modified", fileInfo.ModTime().Round(time.Second).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
		if imsStr := req.Header.Get("If-Modified-Since"); imsStr != "" {
			ims, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", imsStr)
			if err == nil && !fileInfo.ModTime().Round(time.Second).After(ims) {
				ims = ims.Local()
				logger.Debugf("Received `If-Modified-Since` %s is after %s, return 304", imsStr, fileInfo.ModTime())
				w.WriteHeader(http.StatusNotModified) // Not Modified
				return
			} else if err == nil {
				logger.Debugf("Received `If-Modified-Since` %s is before %s", imsStr, fileInfo.ModTime())
			}
		}
	}
	_, err = io.Copy(w, file)
	if err != nil {
		handleFileError(err, path, w)
		return
	}
}

// handleRequestForDirectory handles the request towards the local directory with the given path
func handleRequestForDirectory(w http.ResponseWriter, req *http.Request, path string, fileInfo os.FileInfo) {
	if !enableList {
		indexPath := filepath.Join(path, "index.html")
		fileInfo, err := os.Stat(indexPath)
		if err != nil {
			handleFileError(err, indexPath, w)
			return
		}
		handleRequestForFile(w, req, indexPath, fileInfo)
		return
	}

	// Directory listing is enabled. List the directory
	d, err := os.Open(path)
	if err != nil {
		handleFileError(err, path, w)
		return
	}
	defer d.Close()

	// TODO List the directory using default template
}

func handleFileError(err error, path string, w http.ResponseWriter) {
	if os.IsNotExist(err) {
		logger.Debugf("%s not found.\n", path)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found")
		return
	}
	logger.Warnf("Exception occurred when opening %s: %s\n", path, err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, err)
}
