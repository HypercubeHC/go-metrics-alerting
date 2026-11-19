package handler

import (
	"fmt"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

func GetMainPage(gaugeStorage Repository, counterStorage Repository) string {
	var template string
	template += "<html><head><title>Metrics</title></head><body>"
	template += "<h1>Metrics List</h1>"

	template += "<h2>Gauge metrics</h2><ul>"
	var gaugeMetrics string
	for _, name := range gaugeStorage.Keys() {
		if val, ok := gaugeStorage.Get(name); ok {
			gaugeMetrics += fmt.Sprintf("<li>%s: %s</li>", name, storage.BytesToGauge(val))
		}
	}
	template += gaugeMetrics + "</ul>"

	template += "<h2>Counter metrics</h2><ul>"
	var counterMetrics string
	for _, name := range counterStorage.Keys() {
		if val, ok := counterStorage.Get(name); ok {
			counterMetrics += fmt.Sprintf("<li>%s: %s</li>", name, storage.BytesToCounter(val))
		}
	}
	template += counterMetrics + "</ul>"
	template += "</body></html>"
	return template
}
