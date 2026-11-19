package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/router"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

const testDefaultHeader = "text/plain; charset=utf-8"

type HandlerTestSuite struct {
	suite.Suite
	router http.Handler
}

func TestMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) SetupTest() {
	s.router = router.NewRouter(storage.New[string, []byte](), storage.New[string, []byte]())
}

func doRequest(t *testing.T, handler http.Handler, method, url string) *httptest.ResponseRecorder {
	t.Helper()
	request, err := http.NewRequest(method, url, nil)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func (s *HandlerTestSuite) TestGaugePostHandler() {
	t := s.T()
	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{
			name:       "Gauge correct value",
			url:        "/update/gauge/test/0.8",
			statusCode: http.StatusOK,
		},
		{
			name:       "Gauge correct negative value",
			url:        "/update/gauge/test/-0.5",
			statusCode: http.StatusOK,
		},
		{
			name:       "Gauge incorrect value",
			url:        "/update/gauge/test/abc",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Gauge empty name",
			url:        "/update/gauge//0.8",
			statusCode: http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := doRequest(t, s.router, http.MethodPost, test.url)
			assert.Equal(t, test.statusCode, response.Code)
			if test.statusCode == http.StatusOK {
				assert.Equal(t, testDefaultHeader, response.Header().Get("Content-Type"))
			}
		})
	}
}

func (s *HandlerTestSuite) TestGetGaugeHandlerNegative() {
	t := s.T()

	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{
			name:       "Gauge unknown metric name",
			url:        "/update/gauge/unknown",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Gauge empty metric name",
			url:        "/update/gauge//",
			statusCode: http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := doRequest(t, s.router, http.MethodGet, test.url)
			assert.Equal(t, test.statusCode, response.Code)
		})
	}
}

func (s *HandlerTestSuite) TestGaugeGetHandler() {
	t := s.T()

	tests := []struct {
		name       string
		postURL    string
		getURL     string
		value      storage.Gauge
		statusCode int
	}{
		{
			name:       "Gauge get metric \"Alloc\" is OK",
			postURL:    "/update/gauge/Alloc/0.8",
			getURL:     "/value/gauge/Alloc",
			value:      0.8,
			statusCode: http.StatusOK,
		},
		{
			name:       "Gauge get metric \"test\" is OK",
			postURL:    "/update/gauge/test/1.1",
			getURL:     "/value/gauge/test",
			value:      1.1,
			statusCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, http.StatusOK, doRequest(t, s.router, http.MethodPost, test.postURL).Code)
			response := doRequest(t, s.router, http.MethodGet, test.getURL)
			assert.Equal(t, test.statusCode, response.Code)
			assert.Equal(t, test.value.String(), response.Body.String())
			assert.Equal(t, testDefaultHeader, response.Header().Get("Content-Type"))
		})
	}
}

func (s *HandlerTestSuite) TestCounterPostHandler() {
	t := s.T()
	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{
			name:       "Counter correct value",
			url:        "/update/counter/test/1",
			statusCode: http.StatusOK,
		},
		{
			name:       "Counter correct negative value",
			url:        "/update/counter/test/-1",
			statusCode: http.StatusOK,
		},
		{
			name:       "Counter incorrect value",
			url:        "/update/counter/test/abc",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Counter empty name",
			url:        "/update/counter//1",
			statusCode: http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := doRequest(t, s.router, http.MethodPost, test.url)
			assert.Equal(t, test.statusCode, response.Code)
			if test.statusCode == http.StatusOK {
				assert.Equal(t, testDefaultHeader, response.Header().Get("Content-Type"))
			}
		})
	}
}

func (s *HandlerTestSuite) TestCounterPostHandlerAccumulates() {
	t := s.T()
	t.Run("Counter metric PollCounter sum 1 and 10 equal 11", func(t *testing.T) {
		require.Equal(t, http.StatusOK, doRequest(t, s.router, http.MethodPost, "/update/counter/PollCount/1").Code)
		require.Equal(t, http.StatusOK, doRequest(t, s.router, http.MethodPost, "/update/counter/PollCount/10").Code)

		response := doRequest(t, s.router, http.MethodGet, "/value/counter/PollCount")
		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "11", response.Body.String())
	})

}

func (s *HandlerTestSuite) TestGetCounterHandlerNegative() {
	t := s.T()

	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{
			name:       "Counter unknown metric name",
			url:        "/update/counter/unknown",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Counter empty metric name",
			url:        "/update/counter//",
			statusCode: http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := doRequest(t, s.router, http.MethodGet, test.url)
			assert.Equal(t, test.statusCode, response.Code)
		})
	}
}

func (s *HandlerTestSuite) TestCounterGetHandler() {
	t := s.T()

	tests := []struct {
		name       string
		postURL    string
		getURL     string
		value      storage.Counter
		statusCode int
	}{
		{
			name:       "Counter get metric \"PollCount\" is OK",
			postURL:    "/update/counter/PollCount/1",
			getURL:     "/value/counter/PollCount",
			value:      1,
			statusCode: http.StatusOK,
		},
		{
			name:       "Counter get metric \"Test\" is OK",
			postURL:    "/update/counter/Test/10",
			getURL:     "/value/counter/Test",
			value:      10,
			statusCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, http.StatusOK, doRequest(t, s.router, http.MethodPost, test.postURL).Code)
			response := doRequest(t, s.router, http.MethodGet, test.getURL)
			assert.Equal(t, test.statusCode, response.Code)
			assert.Equal(t, test.value.String(), response.Body.String())
			assert.Equal(t, testDefaultHeader, response.Header().Get("Content-Type"))
		})
	}
}

func (s *HandlerTestSuite) TestFailPostHandler() {
	t := s.T()
	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{
			name:       "Post \"/\" - handler return Not Found status",
			url:        "/",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Post \"/update/\" - handler return bad request status",
			url:        "/update/",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Post \"/value/\" - handler return bad request status",
			url:        "/value/",
			statusCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.statusCode, doRequest(t, s.router, http.MethodPost, test.url).Code)
	}
}

// TODO: добавить тесты хэндлеров UpdatePostHandler, ValuePostHandler
