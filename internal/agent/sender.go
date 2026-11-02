package agent

import (
	"fmt"

	"resty.dev/v3"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

var (
	defaultHeader = "text/plain"
	url           = "http://localhost:8080"
)

func sendPost(client *resty.Client, url string) {
	_, err := client.R().
		SetHeader("Content-Type", defaultHeader).
		Post(url)
	if err != nil {
		return
	}
}

func sendGauge(gaugeStorage *storage.MemStorage[string, gaugeValue], client *resty.Client) {
	for _, name := range gaugeStorage.Keys() {
		if value, ok := gaugeStorage.Get(name); ok {
			url := fmt.Sprintf("%s/update/gauge/%s/%v", url, name, value)
			sendPost(client, url)
		}
	}
}

func sendCounter(counterStorage *storage.MemStorage[string, counterValue], client *resty.Client) {
	for _, name := range counterStorage.Keys() {
		if value, ok := counterStorage.Get(name); ok {
			url := fmt.Sprintf("%s/update/counter/%s/%d", url, name, value)
			sendPost(client, url)
		}
	}
}
