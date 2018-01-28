package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	mime "github.com/mr-dai/http-server/mime"
	"net/url"
)

// handleRequest handles all incoming HTTP requests.
func handleRequest(w http.ResponseWriter, req *http.Request) {
	wrapped := Wrap(&w)
	handleRequestWrapped(&wrapped, req)
	logger.Infof("%s %s %s %d %d", req.Method, req.RequestURI, req.Proto, wrapped.code, wrapped.length)
}

func handleRequestWrapped(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" { // Method not allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path, _ := url.QueryUnescape(req.RequestURI)
	path = filepath.Join(dir, path)
	logger.Tracef("Checking %s", path)
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) { // Append `.html` suffix and try again
		htmlPath := path + ".html"
		fileInfo, err = os.Stat(htmlPath)
	}
	if err != nil {
		handleFileError(err, path, w)
		return
	}

	if fileInfo.IsDir() {
		handleRequestForDirectory(w, req, path, fileInfo)
	} else {
		handleRequestForFile(w, req, path, fileInfo)
	}
}

// handleRequestForFile handles the request towards the given local file.
func handleRequestForFile(w http.ResponseWriter, req *http.Request, path string, fileInfo os.FileInfo) {
	file, err := os.Open(path)
	if err != nil {
		handleFileError(err, path, w)
		return
	}
	defer file.Close()

	w.Header().Add("Server", "mr-dai/http-server/"+version)
	w.Header().Add("Content-Type", mime.GetMimeByFilename(file.Name()))

	if enableCache { // Handle HTTP cache headers
		w.Header().Add("Last-Modified", fileInfo.ModTime().Round(time.Second).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
		if imsStr := req.Header.Get("If-Modified-Since"); imsStr != "" {
			ims, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", imsStr)
			if err == nil && !fileInfo.ModTime().Round(time.Second).After(ims) {
				ims = ims.Local()
				logger.Debugf("Received `If-Modified-Since` %s is after %s, return 304.", imsStr, fileInfo.ModTime())
				w.WriteHeader(http.StatusNotModified) // Not Modified
				return
			} else if err == nil {
				logger.Debugf("Received `If-Modified-Since` %s is before %s.", imsStr, fileInfo.ModTime())
			}
		}
		if maxAge > 0 {
			w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
		} else {
			w.Header().Add("Cache-Control", "no-cache")
		}
	} else { // Disallow clients to cache the response
		w.Header().Add("Cache-Control", "no-store")
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
	names, err := d.Readdirnames(0)
	if err != nil {
		handleFileError(err, path, w)
		return
	}

	// List the directory using default template
	listDir(req.RequestURI, names, w)
}

// handleFileError handles the given error and returns an appropriate HTTP response.
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
