package agent

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

type Repository interface {
	Get(name string) ([]byte, bool)
	Set(name string, val []byte)
	Keys() []string
	Iter(ctx context.Context, fn func(name string, val []byte) error) error
	Size() int
	Delete(name string)
}

func getMetrics(ctx context.Context, timeInterval time.Duration, storages map[string]*metricStorage) {

	timer := time.NewTicker(timeInterval)
	defer timer.Stop()

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			for name := range storages {
				switch name {
				case gauge:
					getGaugeMetrics(storages[name])
				case counter:
					getCounterMetrics(storages[name])
				default:
				}
			}
		}
	}
}

func getGaugeMetrics(gaugeStorage *metricStorage) {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	gaugeStorage.mx.Lock()
	defer gaugeStorage.mx.Unlock()

	gaugeStorage.Set("Alloc", storage.GaugeToBytes(storage.Gauge(memStats.Alloc)))
	gaugeStorage.Set("BuckHashSys", storage.GaugeToBytes(storage.Gauge(memStats.BuckHashSys)))
	gaugeStorage.Set("Frees", storage.GaugeToBytes(storage.Gauge(memStats.Frees)))
	gaugeStorage.Set("GCCPUFraction", storage.GaugeToBytes(storage.Gauge(memStats.GCCPUFraction)))
	gaugeStorage.Set("GCSys", storage.GaugeToBytes(storage.Gauge(memStats.GCSys)))
	gaugeStorage.Set("HeapAlloc", storage.GaugeToBytes(storage.Gauge(memStats.HeapAlloc)))
	gaugeStorage.Set("HeapIdle", storage.GaugeToBytes(storage.Gauge(memStats.HeapIdle)))
	gaugeStorage.Set("HeapInuse", storage.GaugeToBytes(storage.Gauge(memStats.HeapInuse)))
	gaugeStorage.Set("HeapObjects", storage.GaugeToBytes(storage.Gauge(memStats.HeapObjects)))
	gaugeStorage.Set("HeapReleased", storage.GaugeToBytes(storage.Gauge(memStats.HeapReleased)))
	gaugeStorage.Set("HeapSys", storage.GaugeToBytes(storage.Gauge(memStats.HeapSys)))
	gaugeStorage.Set("LastGC", storage.GaugeToBytes(storage.Gauge(memStats.LastGC)))
	gaugeStorage.Set("Lookups", storage.GaugeToBytes(storage.Gauge(memStats.Lookups)))
	gaugeStorage.Set("MCacheInuse", storage.GaugeToBytes(storage.Gauge(memStats.MCacheInuse)))
	gaugeStorage.Set("MCacheSys", storage.GaugeToBytes(storage.Gauge(memStats.MCacheSys)))
	gaugeStorage.Set("MSpanInuse", storage.GaugeToBytes(storage.Gauge(memStats.MSpanInuse)))
	gaugeStorage.Set("MSpanSys", storage.GaugeToBytes(storage.Gauge(memStats.MSpanSys)))
	gaugeStorage.Set("Mallocs", storage.GaugeToBytes(storage.Gauge(memStats.Mallocs)))
	gaugeStorage.Set("NextGC", storage.GaugeToBytes(storage.Gauge(memStats.NextGC)))
	gaugeStorage.Set("NumForcedGC", storage.GaugeToBytes(storage.Gauge(memStats.NumForcedGC)))
	gaugeStorage.Set("NumGC", storage.GaugeToBytes(storage.Gauge(memStats.NumGC)))
	gaugeStorage.Set("OtherSys", storage.GaugeToBytes(storage.Gauge(memStats.OtherSys)))
	gaugeStorage.Set("PauseTotalNs", storage.GaugeToBytes(storage.Gauge(memStats.PauseTotalNs)))
	gaugeStorage.Set("StackInuse", storage.GaugeToBytes(storage.Gauge(memStats.StackInuse)))
	gaugeStorage.Set("StackSys", storage.GaugeToBytes(storage.Gauge(memStats.StackSys)))
	gaugeStorage.Set("Sys", storage.GaugeToBytes(storage.Gauge(memStats.Sys)))
	gaugeStorage.Set("TotalAlloc", storage.GaugeToBytes(storage.Gauge(memStats.TotalAlloc)))
	gaugeStorage.Set("RandomValue", storage.GaugeToBytes(storage.Gauge(rand.Float64())))
}

func getCounterMetrics(conuterStorage *metricStorage) {
	conuterStorage.mx.Lock()
	defer conuterStorage.mx.Unlock()

	counter, ok := conuterStorage.Get("PollCount")
	if !ok {
		conuterStorage.Set("PollCount", storage.CounterToBytes(storage.Counter(1)))
		return
	}
	counterUpdate := storage.BytesToCounter(counter)
	counterUpdate++
	conuterStorage.Set("PollCount", storage.CounterToBytes(counterUpdate))
}
