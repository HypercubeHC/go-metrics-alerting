package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/router"
)

var testName string = "test"
var headerName string = "Content-Type"
var testDefaultHeader http.Header = http.Header{headerName: []string{"text/plain; charset=utf-8"}}

func TestServer(t *testing.T) {
	server := httptest.NewServer(router.NewRouter())
	defer server.Close()

	gaugeTestVal := float64(99.9)
	gaugeTestValStr := strconv.FormatFloat(gaugeTestVal, 'g', -1, 64)
	counterTestVal := int64(2)
	counterTestValStr := strconv.FormatInt(counterTestVal, 10)

	tests := []struct {
		name       string
		url        string
		method     string
		header     http.Header
		statusCode int
	}{
		{
			name:       "Positive. Gauge. Happy path",
			url:        server.URL + "/update/gauge/" + testName + "/" + gaugeTestValStr,
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusOK,
		},
		{
			name:       "Negative. Gauge. Incorrent value",
			url:        server.URL + "/update/gauge/" + testName + "/" + "abc",
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Negative. Gauge. Empty gauge name",
			url:        server.URL + "/update/gauge/" + "/" + gaugeTestValStr,
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Negative. Gauge. Incorrect header",
			url:        server.URL + "/update/gauge/" + testName + "/" + gaugeTestValStr,
			method:     http.MethodPost,
			header:     http.Header{headerName: []string{"application/json"}},
			statusCode: http.StatusOK,
		},
		{
			name:       "Positive. Counter. Happy path",
			url:        server.URL + "/update/counter/" + testName + "/" + counterTestValStr,
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusOK,
		},
		{
			name:       "Negative. Counter. Incorrent value",
			url:        server.URL + "/update/counter/" + testName + "/" + "09-ff",
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Negative. Counter. Empty counter name",
			url:        server.URL + "/update/counter/" + "/" + counterTestValStr,
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Negative. Gauge. Incorrect header",
			url:        server.URL + "/update/gauge/" + testName + "/" + gaugeTestValStr,
			method:     http.MethodPost,
			header:     http.Header{headerName: []string{"application/json"}},
			statusCode: http.StatusOK,
		},
		{
			name:       "Negative. Wrong metric type",
			url:        server.URL + "/update/another/" + testName + "/" + gaugeTestValStr,
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Negative. Empty metric name",
			url:        server.URL + "/update/gauge//" + gaugeTestValStr,
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Negative. URL without metric type, name and value",
			url:        server.URL + "/update/",
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Negative. Empty URL",
			url:        server.URL + "/",
			method:     http.MethodPost,
			header:     testDefaultHeader,
			statusCode: http.StatusNotFound,
		},
	}

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatalf("Failed request: %v", err)
			}
			for key, values := range test.header {
				for _, value := range values {
					request.Header.Add(key, value)
				}
			}
			response, err := client.Do(request)
			if err != nil {
				t.Fatalf("Failed response: %v", err)
			}
			defer response.Body.Close()
			expectedContentType := request.Header.Get(headerName)
			actualContentType := response.Header.Get(headerName)

			assert.Equal(t, test.statusCode, response.StatusCode)

			if actualContentType != expectedContentType {
				assert.NotEqual(t, expectedContentType, actualContentType)
			} else {
				assert.Equal(t, expectedContentType, actualContentType)
			}

		})
	}
}
