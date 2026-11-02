package template

import (
	"fmt"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

func GetMainPage(gaugeStorage *storage.MemStorage[string, float64], counterStorage *storage.MemStorage[string, int64]) string {
	var template string
	template += "<html><head><title>Metrics</title></head><body>"
	template += "<h1>Metrics List</h1>"

	template += "<h2>Gauge metrics</h2><ul>"
	var gaugeMetrics string
	for _, name := range gaugeStorage.Keys() {
		if val, ok := gaugeStorage.Get(name); ok {
			gaugeMetrics += fmt.Sprintf("<li>%s: %v</li>", name, val)
		}
	}
	template += gaugeMetrics + "</ul>"

	template += "<h2>Counter metrics</h2><ul>"
	var counterMetrics string
	for _, name := range counterStorage.Keys() {
		if val, ok := counterStorage.Get(name); ok {
			counterMetrics += fmt.Sprintf("<li>%s: %d</li>", name, val)
		}
	}
	template += counterMetrics + "</ul>"
	template += "</body></html>"
	return template
}
