# http-server

Static HTTP server in Golang.

This is just one of my practice projects. Don't expect it to be ready for production usage.

## Usage

Use `go build` or `go install` to build or install the program.

## Future Tasks

- [x] Support for common file types(`text/html`, `text/json`, ...)
- [x] Full support for HTTP cache (`If-Modified-Since`, `Last-Modified`, `Cache-Control`)
- [ ] List directory using default HTML template
- [ ] Write `Makefile`, make sure it can be easilly installed
- [ ] Add test cases
- [ ] Add doc comments
- [ ] Integrate with Travis CI and Codecov.io.
- [ ] Release version 0.1
- [ ] Support for GZip
- [ ] Support for HTTPS
- [ ] Support for secure file upload with token-based authentication
- [ ] Configuration for controlling access from local network and Internet (default to `localhost` only)
- [ ] Release version 0.2
- [ ] Use [fsnotify](https://github.com/fsnotify/fsnotify/) to implement a more efficient support for HTTP cache

## Side Tasks

- [ ] Color output for `log`.
- [ ] Implement native in-memory cache
- [ ] Implement `CacheManager` for cache using other source, e.g. Redis.
