package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (resWriter *responseWriter) WriteHeader(code int) {
	resWriter.statusCode = code
	resWriter.ResponseWriter.WriteHeader(code)
}

func (resWriter *responseWriter) Write(bytesToWrite []byte) (int, error) {
	if resWriter.statusCode == 0 {
		resWriter.statusCode = http.StatusOK
	}
	size, err := resWriter.ResponseWriter.Write(bytesToWrite)
	resWriter.responseSize += size
	return size, err
}

func LogMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(logMiddleware http.Handler) http.Handler {
		return http.HandlerFunc(func(resWriter http.ResponseWriter, request *http.Request) {
			start := time.Now()
			midResWriter := &responseWriter{ResponseWriter: resWriter}
			logMiddleware.ServeHTTP(midResWriter, request)

			logger.Infow("*INFO* request",
				"uri", request.RequestURI,
				"method", request.Method,
				"timeExecution", time.Since(start),
				"statusCode", midResWriter.statusCode,
				"responseSize", midResWriter.responseSize,
			)
		})
	}
}
