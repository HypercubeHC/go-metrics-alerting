package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

type MetricsTestSuite struct {
	suite.Suite
	storages map[string]*storage.MemStorage[string, []byte]
}

func TestMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsTestSuite))
}

func (s *MetricsTestSuite) SetupTest() {
	s.storages = createStorages()
}

func (s *MetricsTestSuite) TestCreateStorages() {
	t := s.T()
	assert.Len(t, s.storages, 2)
	assert.Contains(t, s.storages, gauge)
	assert.Contains(t, s.storages, counter)
}

func (s *MetricsTestSuite) TestGetGaugeMetrics() {
	t := s.T()
	getGaugeMetrics(s.storages[gauge])

	keys := s.storages[gauge].Keys()
	assert.Len(t, keys, 28)

	for _, name := range []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	} {
		_, ok := s.storages[gauge].Get(name)
		assert.True(t, ok, "Gauge metric: \"%s\" - should be collected, not present", name)
	}
}

func (s *MetricsTestSuite) TestGetCounterMetrics() {
	t := s.T()
	getCounterMetrics(s.storages[counter])
	keys := s.storages[counter].Keys()
	assert.Len(t, keys, 1)

	for _, name := range []string{"PollCount"} {
		_, ok := s.storages[counter].Get(name)
		assert.True(t, ok, "Counter metric: \"%s\" - should be collected, not present", name)
	}
}

func (s *MetricsTestSuite) TestGaugeRandomValue() {
	t := s.T()
	getGaugeMetrics(s.storages[gauge])
	randomValue, ok := s.storages[gauge].Get("RandomValue")
	require.True(t, ok)

	value := storage.BytesToGauge(randomValue)
	assert.GreaterOrEqual(t, value, storage.Gauge(0))
	assert.Less(t, value, storage.Gauge(1))
}

func (s *MetricsTestSuite) TestCounterMetricsAccumulates() {
	t := s.T()

	getCounterMetrics(s.storages[counter])
	getCounterMetrics(s.storages[counter])

	counterValue, ok := s.storages[counter].Get("PollCount")
	require.True(t, ok)
	assert.Equal(t, storage.Counter(2), storage.BytesToCounter(counterValue))

	getCounterMetrics(s.storages[counter])
	newCounterValue, _ := s.storages[counter].Get("PollCount")
	assert.Equal(t, storage.Counter(3), storage.BytesToCounter(newCounterValue))
}
