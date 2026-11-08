package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
	"resty.dev/v3"
)

var (
	defaultHeader = "text/plain"
)

func sendMetrics(ctx context.Context, url string, timeInterval time.Duration, client *resty.Client,
	storages map[string]*metricStorage) {

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
					sendGauge(ctx, url, storages[name], client)
				case counter:
					sendCounter(ctx, url, storages[name], client)
				default:
				}
			}
		}
	}
}

func sendGauge(ctx context.Context, url string, gaugeStorage *metricStorage, client *resty.Client) {
	gaugeStorage.mx.Lock()
	defer gaugeStorage.mx.Unlock()

	err := gaugeStorage.Iter(ctx, func(name string, val []byte) error {
		urlPost := "http://" + url + "/update/gauge/" + name + "/" + storage.BytesToGauge(val).String()
		response, err := client.R().
			SetHeader("Content-Type", defaultHeader).
			Post(urlPost)
		if err != nil {
			return err
		}
		return response.Body.Close()
	})
	if err != nil {
		fmt.Printf("Failed to send gauge metric; %s\n", err)
	}
}

func sendCounter(ctx context.Context, url string, counterStorage *metricStorage, client *resty.Client) {
	counterStorage.mx.Lock()
	defer counterStorage.mx.Unlock()

	err := counterStorage.Iter(ctx, func(name string, val []byte) error {
		urlPost := "http://" + url + "/update/counter/" + name + "/" + storage.BytesToCounter(val).String()
		response, err := client.R().
			SetHeader("Content-Type", defaultHeader).
			Post(urlPost)
		if err != nil {
			return err
		}
		return response.Body.Close()
	})
	if err != nil {
		fmt.Printf("Failed to send counter metric; %s\n", err)
	}
}
