package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	log "github.com/mr-dai/http-server/log"
)

const version = "0.1"

var logger *log.Logger

var enableCache bool
var maxAge int
var enableList bool
var port int
var dir string

var debug bool
var trace bool

func setupFlag() {
	flag.StringVar(&dir, "dir", "", "path of directory to server, default to be the current working directory")
	flag.IntVar(&port, "port", 8080, "specify the port number to listen on")

	flag.BoolVar(&enableCache, "cache", false, "enable HTTP cache support")
	flag.IntVar(&maxAge, "maxAge", -1, "`max-age` field for `Cache-Control` header. Used only when cache support is enabled")

	flag.BoolVar(&enableList, "list", false, "enable listing on directory when index.html is missed")
	flag.BoolVar(&debug, "debug", false, "enable DEBUG level log output")
	flag.BoolVar(&trace, "trace", false, "enable TRACE level log output")
	flag.Parse()

	if dir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			logger.Fatal(err.Error())
		}
		dir = pwd
	}
}

func setupLogging() {
	if trace {
		logger = log.NewLogger(log.TRACE, "http-server", "")
	} else if debug {
		logger = log.NewLogger(log.DEBUG, "http-server", "")
	} else {
		logger = log.NewLogger(log.INFO, "http-server", "")
	}
}

func main() {
	// Initialize
	setupFlag()
	setupLogging()

	logger.Debugf("Received configuration { port = `%d`, dir = `%s`, enableList = `%t`, enableCache = `%t` }",
		port, dir, enableList, enableCache)
	logger.Infof("Server is listening on `%d`...", port)

	// Setup HTTP server
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(fmt.Sprintf(":%d", port), http.DefaultServeMux)

	// TODO Wait on interrupt signal

	logger.Info("Received interrupt signal, exiting")
}
