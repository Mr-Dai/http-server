# static-server

Static HTTP server in Golang.

This is just one of my practice projects. Don't expect it to be ready for production usage.

## Usage

Use `go build` or `go install` to build or install the program.

## Future Tasks

- [x] Support for common file types(`text/html`, `text/json`, ...)
- [ ] Support for HTTP cache(`ETag`, `Cache-Control`, 
      Refer to [this](https://developers.google.com/web/fundamentals/performance/optimizing-content-efficiency/http-caching))
- [ ] Implement full-fledged log framework. (mimic [log4j](http://logging.apache.org/log4j/2.x/index.html))
- [ ] List directory using default HTML template
- [ ] Write `Makefile`, make sure it can be easilly installed
- [ ] Add test cases
- [ ] Integrate with Travis CI and Codecov.io.
- [ ] Release version 0.1
- [ ] Support for GZip
- [ ] Support for HTTPS
- [ ] Support for secure file upload with token-based authentication
- [ ] Configuration for controlling access from local network and Internet (default to `localhost` only)
- [ ] Release version 0.2
- [ ] Use [fsnotify](https://github.com/fsnotify/fsnotify/) to implement a more efficient support for HTTP cache

## Side Tasks

- [ ] Write benchmark tests to compare the performance of my log framework and other popular log framework,
      e.g. [zap](https://github.com/uber-go/zap)
- [ ] Refactor the log framework to use goroutines and channels to implement Producer-Consumer pattern.
