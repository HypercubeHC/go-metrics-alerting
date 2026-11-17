package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

const (
	jsonHeader = "application/json"
	htmlHeader = "text/html"
	gzipHeader = "gzip"
)

type compressResponseWriter struct {
	http.ResponseWriter
	compressWriter *gzip.Writer
	compress       bool
	started        bool
}

func (rw *compressResponseWriter) start() {
	if rw.started {
		return
	}
	rw.started = true
	contentType := rw.Header().Get("Content-Type")
	if strings.HasPrefix(contentType, jsonHeader) || strings.HasPrefix(contentType, htmlHeader) {
		rw.Header().Set("Content-Encoding", gzipHeader)
		rw.compressWriter = gzip.NewWriter(rw.ResponseWriter)
		rw.compress = true
	}
}

func (rw *compressResponseWriter) Close() {
	if rw.compressWriter != nil {
		rw.compressWriter.Close()
	}
}

func (rw *compressResponseWriter) WriteHeader(code int) {
	rw.start()
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *compressResponseWriter) Write(bytesToWrtire []byte) (int, error) {
	rw.start()
	if rw.compress {
		return rw.compressWriter.Write(bytesToWrtire)
	}
	return rw.ResponseWriter.Write(bytesToWrtire)
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Content-Encoding") == gzipHeader {
			compressReader, err := gzip.NewReader(request.Body)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			defer compressReader.Close()
			request.Body = compressReader
		}

		if !strings.Contains(request.Header.Get("Accept-Encoding"), gzipHeader) {
			next.ServeHTTP(rw, request)
			return
		}

		wrapped := &compressResponseWriter{ResponseWriter: rw}
		defer wrapped.Close()
		next.ServeHTTP(wrapped, request)
	})
}
