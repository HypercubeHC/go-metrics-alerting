package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

var (
	url            = "http://localhost:8080"
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

type (
	gaugeValue   float64
	counterValue int64
)

func getGaugeMetrics(gaugeStorage *storage.MemStorage[string, gaugeValue]) {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	gaugeStorage.Set("Alloc", gaugeValue(memStats.Alloc))
	gaugeStorage.Set("BuckHashSys", gaugeValue(memStats.BuckHashSys))
	gaugeStorage.Set("Frees", gaugeValue(memStats.Frees))
	gaugeStorage.Set("GCCPUFraction", gaugeValue(memStats.GCCPUFraction))
	gaugeStorage.Set("GCSys", gaugeValue(memStats.GCSys))
	gaugeStorage.Set("HeapAlloc", gaugeValue(memStats.HeapAlloc))
	gaugeStorage.Set("HeapIdle", gaugeValue(memStats.HeapIdle))
	gaugeStorage.Set("HeapInuse", gaugeValue(memStats.HeapInuse))
	gaugeStorage.Set("HeapObjects", gaugeValue(memStats.HeapObjects))
	gaugeStorage.Set("HeapReleased", gaugeValue(memStats.HeapReleased))
	gaugeStorage.Set("HeapSys", gaugeValue(memStats.HeapSys))
	gaugeStorage.Set("LastGC", gaugeValue(memStats.LastGC))
	gaugeStorage.Set("Lookups", gaugeValue(memStats.Lookups))
	gaugeStorage.Set("MCacheInuse", gaugeValue(memStats.MCacheInuse))
	gaugeStorage.Set("MCacheSys", gaugeValue(memStats.MCacheSys))
	gaugeStorage.Set("MSpanInuse", gaugeValue(memStats.MSpanInuse))
	gaugeStorage.Set("MSpanSys", gaugeValue(memStats.MSpanSys))
	gaugeStorage.Set("Mallocs", gaugeValue(memStats.Mallocs))
	gaugeStorage.Set("NextGC", gaugeValue(memStats.NextGC))
	gaugeStorage.Set("NumForcedGC", gaugeValue(memStats.NumForcedGC))
	gaugeStorage.Set("NumGC", gaugeValue(memStats.NumGC))
	gaugeStorage.Set("OtherSys", gaugeValue(memStats.OtherSys))
	gaugeStorage.Set("PauseTotalNs", gaugeValue(memStats.PauseTotalNs))
	gaugeStorage.Set("StackInuse", gaugeValue(memStats.StackInuse))
	gaugeStorage.Set("StackSys", gaugeValue(memStats.StackSys))
	gaugeStorage.Set("Sys", gaugeValue(memStats.Sys))
	gaugeStorage.Set("TotalAlloc", gaugeValue(memStats.TotalAlloc))
	gaugeStorage.Set("RandomValue", gaugeValue(rand.Float64()))
}

func getCounterMetrics(counterStorage *storage.MemStorage[string, counterValue]) {
	pollCount, ok := counterStorage.Get("PollCount")
	if !ok {
		counterStorage.Set("PollCount", 1)
		return
	}
	pollCount++
	counterStorage.Set("PollCount", pollCount)
}

func sendPost(client *http.Client, url string) {
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "text/plain")

	response, err := client.Do(request)
	if err != nil {
		return
	}
	response.Body.Close()
}

func main() {
	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	gaugeStorage := storage.New[string, gaugeValue]()
	counterStorage := storage.New[string, counterValue]()

	client := &http.Client{}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-pollTicker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			getGaugeMetrics(gaugeStorage)
			getCounterMetrics(counterStorage)

		case <-reportTicker.C:
			for _, name := range gaugeStorage.Keys() {
				if value, ok := gaugeStorage.Get(name); ok {
					url := fmt.Sprintf("%s/update/gauge/%s/%v", url, name, value)
					sendPost(client, url)
				}
			}

			for _, name := range counterStorage.Keys() {
				if value, ok := counterStorage.Get(name); ok {
					url := fmt.Sprintf("%s/update/counter/%s/%d", url, name, value)
					sendPost(client, url)
				}
			}
		case sig := <-stop:
			fmt.Printf("Agent got signal: %v.\n", sig)
			return
		}
	}
}
