package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogPostGaugeRequest(t *testing.T) {
	tests := []struct {
		name                 string
		handler              http.HandlerFunc
		logLevel             zapcore.Level
		method               string
		uri                  string
		expectedStatusCode   int64
		expectedResponseSize int64
	}{
		{
			name: "Log POST gauge with writed body",
			handler: func(rw http.ResponseWriter, r *http.Request) {
				_, _ = rw.Write([]byte("ok"))
			},
			logLevel:             zap.InfoLevel,
			method:               http.MethodPost,
			uri:                  "/update/gauge/test/1.1",
			expectedStatusCode:   http.StatusOK,
			expectedResponseSize: 2,
		},
		{
			name: "Log POST gauge with writed 201 status and body",
			handler: func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusCreated)
				_, _ = rw.Write([]byte("test"))
			},
			logLevel:             zap.InfoLevel,
			method:               http.MethodPost,
			uri:                  "/update/gauge/Alloc/0.8",
			expectedStatusCode:   http.StatusCreated,
			expectedResponseSize: 4,
		},
		{
			name: "Log POST gauge with writed 404 status",
			handler: func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			},
			logLevel:             zap.InfoLevel,
			method:               http.MethodPost,
			uri:                  "/update/gauge/Alloc/1.1",
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseSize: 0,
		},
		{
			name: "Log POST counter with writed body",
			handler: func(rw http.ResponseWriter, r *http.Request) {
				_, _ = rw.Write([]byte("ok"))
			},
			logLevel:             zap.InfoLevel,
			method:               http.MethodPost,
			uri:                  "/update/counter/test/1",
			expectedStatusCode:   http.StatusOK,
			expectedResponseSize: 2,
		},
		{
			name: "Log POST counter with writed 201 status and body",
			handler: func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusCreated)
				_, _ = rw.Write([]byte("testtest"))
			},
			logLevel:             zap.InfoLevel,
			method:               http.MethodPost,
			uri:                  "/update/counter/PollCount/2",
			expectedStatusCode:   http.StatusCreated,
			expectedResponseSize: 8,
		},
		{
			name: "Log POST counter with writed 404 status",
			handler: func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			},
			logLevel:             zap.InfoLevel,
			method:               http.MethodPost,
			uri:                  "/update/counter/PollCount/3",
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseSize: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			core, observedLogs := observer.New(test.logLevel)
			testLogger := zap.New(core).Sugar()
			request := httptest.NewRequest(test.method, test.uri, nil)
			testResWriter := httptest.NewRecorder()

			LogMiddleware(testLogger)(http.HandlerFunc(test.handler)).ServeHTTP(testResWriter, request)

			require.Equal(t, 1, observedLogs.Len())
			entry := observedLogs.All()[0]
			assert.Equal(t, test.logLevel, entry.Level)
			fields := entry.ContextMap()
			assert.Equal(t, test.uri, fields["uri"])
			assert.Equal(t, test.method, fields["method"])
			assert.Equal(t, test.expectedStatusCode, fields["statusCode"])
			assert.Equal(t, test.expectedResponseSize, fields["responseSize"])
		})
	}
}
