package agent

import (
	"fmt"
	"net/http"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

var (
	defaultHeader = "text/plain"
	url           = "http://localhost:8080"
)

func sendPost(client *http.Client, url string) {
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", defaultHeader)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	response.Body.Close()
}

func sendGauge(gaugeStorage *storage.MemStorage[string, gaugeValue], client *http.Client) {
	for _, name := range gaugeStorage.Keys() {
		if value, ok := gaugeStorage.Get(name); ok {
			url := fmt.Sprintf("%s/update/gauge/%s/%v", url, name, value)
			sendPost(client, url)
		}
	}
}

func sendCounter(counterStorage *storage.MemStorage[string, counterValue], client *http.Client) {
	for _, name := range counterStorage.Keys() {
		if value, ok := counterStorage.Get(name); ok {
			url := fmt.Sprintf("%s/update/counter/%s/%d", url, name, value)
			sendPost(client, url)
		}
	}
}
