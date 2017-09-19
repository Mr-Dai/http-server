package main

import (
	"net/http"
)

// ResponseWriterWrapper wraps an `http.ResponseWriter`. ResponseWriterWrapper is
// responsible for recording response length when its `Write` method is invoked.
type ResponseWriterWrapper struct {
	wrapped *http.ResponseWriter
	length  int
	code    int
}

// Wrap wraps the given `http.ResponseWriter` in a `ResponseWriterWrapper`.
func Wrap(w *http.ResponseWriter) ResponseWriterWrapper {
	return ResponseWriterWrapper{wrapped: w, length: 0, code: 0}
}

// Header returns the `http.Header` of the underlying `http.RespsoneWriter`
func (w ResponseWriterWrapper) Header() http.Header {
	return (*w.wrapped).Header()
}

func (w *ResponseWriterWrapper) Write(bytes []byte) (int, error) {
	len, err := (*w.wrapped).Write(bytes)
	if err == nil {
		w.length += len
		if w.code == 0 {
			w.code = http.StatusOK
		}
	}
	return len, err
}

func (w *ResponseWriterWrapper) WriteHeader(code int) {
	w.code = code
	(*w.wrapped).WriteHeader(code)
}
