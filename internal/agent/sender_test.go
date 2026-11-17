package agent

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"resty.dev/v3"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/model"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

const jsonTestHeader = "application/json"

type SenderTestSuite struct {
	suite.Suite
	client *resty.Client
}

func TestSenderTestSuite(t *testing.T) {
	suite.Run(t, new(SenderTestSuite))
}

func (s *SenderTestSuite) SetupTest() {
	s.client = resty.New()
}

func createTestServer(collected *[]model.Metric, contentType *string, mx *sync.Mutex) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		if r.Header.Get("Content-Encoding") == "gzip" {
			compressReader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer compressReader.Close()
			body = compressReader
		}
		var metric model.Metric
		if err := json.NewDecoder(body).Decode(&metric); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		mx.Lock()
		*collected = append(*collected, metric)
		*contentType = r.Header.Get("Content-Type")
		mx.Unlock()
		w.Header().Set("Content-Type", jsonTestHeader)
		w.WriteHeader(http.StatusOK)
	}))
}

func (s *SenderTestSuite) TestSendCorrectGauge() {
	t := s.T()
	var collected []model.Metric
	var contentType string
	var mx sync.Mutex
	server := createTestServer(&collected, &contentType, &mx)
	defer server.Close()

	tests := []struct {
		name        string
		metricName  string
		metricValue storage.Gauge
	}{
		{
			name:        "Alloc correct value",
			metricName:  "Alloc",
			metricValue: 0.9,
		},
		{
			name:        "BuckHashSys correct value",
			metricName:  "BuckHashSys",
			metricValue: 0.3322,
		},
		{
			name:        "Frees correct value",
			metricName:  "Frees",
			metricValue: 0.0011,
		},
		{
			name:        "GCCPUFraction correct value",
			metricName:  "GCCPUFraction",
			metricValue: 0.2213,
		},
		{
			name:        "GCSys correct value",
			metricName:  "GCSys",
			metricValue: 0.2323,
		},
		{
			name:        "HeapAlloc correct value",
			metricName:  "HeapAlloc",
			metricValue: 0.2132,
		},
		{
			name:        "HeapIdle correct value",
			metricName:  "HeapIdle",
			metricValue: 0.1245,
		},
		{
			name:        "HeapInuse correct value",
			metricName:  "HeapInuse",
			metricValue: 0.19938,
		},
		{
			name:        "HeapObjects correct value",
			metricName:  "HeapObjects",
			metricValue: 0.9977,
		},
		{
			name:        "HeapReleased correct value",
			metricName:  "HeapReleased",
			metricValue: 0.2335,
		},
		{
			name:        "HeapSys correct value",
			metricName:  "HeapSys",
			metricValue: 0.688,
		},
		{
			name:        "LastGC correct value",
			metricName:  "LastGC",
			metricValue: 1.3444,
		},
		{
			name:        "Lookups correct value",
			metricName:  "Lookups",
			metricValue: 94.444,
		},
		{
			name:        "MCacheInuse correct value",
			metricName:  "MCacheInuse",
			metricValue: 0.2223,
		},
		{
			name:        "MCacheSys correct value",
			metricName:  "MCacheSys",
			metricValue: 0.1,
		},
		{
			name:        "MSpanInuse correct value",
			metricName:  "MSpanInuse",
			metricValue: 0.2,
		},
		{
			name:        "MSpanSys correct value",
			metricName:  "MSpanSys",
			metricValue: 9.2,
		},
		{
			name:        "Mallocs correct value",
			metricName:  "Mallocs",
			metricValue: 23.33,
		},
		{
			name:        "NextGC correct value",
			metricName:  "NextGC",
			metricValue: 8.34,
		},
		{
			name:        "NumForcedGC correct value",
			metricName:  "NumForcedGC",
			metricValue: 0.4,
		},
		{
			name:        "NumGC correct value",
			metricName:  "NumGC",
			metricValue: 0.111,
		},
		{
			name:        "OtherSys correct value",
			metricName:  "OtherSys",
			metricValue: 0.7,
		},
		{
			name:        "PauseTotalNs correct value",
			metricName:  "PauseTotalNs",
			metricValue: 10.10,
		},
		{
			name:        "StackInuse correct value",
			metricName:  "StackInuse",
			metricValue: 76.33,
		},
		{
			name:        "StackSys correct value",
			metricName:  "StackSys",
			metricValue: 99.1,
		},
		{
			name:        "Sys correct value",
			metricName:  "Sys",
			metricValue: 444.66,
		},
		{
			name:        "TotalAlloc correct value",
			metricName:  "TotalAlloc",
			metricValue: 0.8,
		},
		{
			name:        "RandomValue correct value",
			metricName:  "RandomValue",
			metricValue: 111.11111,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gaugeStorage := &metricStorage{Repository: storage.New[string, []byte]()}
			gaugeStorage.Set(test.metricName, storage.GaugeToBytes(test.metricValue))

			mx.Lock()
			collected = nil
			contentType = ""
			mx.Unlock()

			sendGauge(context.Background(), server.URL, gaugeStorage, s.client)

			value := test.metricValue
			expected := []model.Metric{{ID: test.metricName, MType: gauge, Value: &value}}
			mx.Lock()
			defer mx.Unlock()
			assert.Equal(t, expected, collected)
			assert.Equal(t, jsonTestHeader, contentType)
		})
	}
}

func (s *SenderTestSuite) TestSendCorrectCounter() {
	t := s.T()
	var collected []model.Metric
	var contentType string
	var mx sync.Mutex
	server := createTestServer(&collected, &contentType, &mx)
	defer server.Close()
	tests := []struct {
		name        string
		metricName  string
		metricValue storage.Counter
	}{
		{
			name:        "PollCount correct value",
			metricName:  "PollCount",
			metricValue: 7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counterStorage := &metricStorage{Repository: storage.New[string, []byte]()}
			counterStorage.Set(test.metricName, storage.CounterToBytes(test.metricValue))

			mx.Lock()
			collected = nil
			contentType = ""
			mx.Unlock()

			sendCounter(context.Background(), server.URL, counterStorage, s.client)

			delta := test.metricValue
			expected := []model.Metric{{ID: test.metricName, MType: counter, Delta: &delta}}
			mx.Lock()
			assert.Equal(t, expected, collected)
			assert.Equal(t, jsonTestHeader, contentType)
			mx.Unlock()
		})
	}
}
