package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/model"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
	"resty.dev/v3"
)

var (
	jsonHeader = "application/json"
)

func sendMetrics(ctx context.Context, url string, timeInterval time.Duration, client *resty.Client,
	storages map[string]*metricStorage) {

	timer := time.NewTicker(timeInterval)
	defer timer.Stop()

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

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

func postMetric(url string, metric model.Metric, client *resty.Client) error {
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	response, err := client.R().
		SetHeader("Content-Type", jsonHeader).
		SetBody(body).
		Post(url + "/update")
	if err != nil {
		return err
	}
	return response.Body.Close()
}

func sendGauge(ctx context.Context, url string, gaugeStorage *metricStorage, client *resty.Client) {
	gaugeStorage.mx.Lock()
	defer gaugeStorage.mx.Unlock()

	err := gaugeStorage.Iter(ctx, func(name string, val []byte) error {
		value := storage.BytesToGauge(val)
		gaugeMetric := model.Metric{ID: name, MType: gauge, Value: &value}
		return postMetric(url, gaugeMetric, client)
	})
	if err != nil {
		fmt.Printf("Failed to send gauge metric; %s\n", err)
	}
}

func sendCounter(ctx context.Context, url string, counterStorage *metricStorage, client *resty.Client) {
	counterStorage.mx.Lock()
	defer counterStorage.mx.Unlock()

	err := counterStorage.Iter(ctx, func(name string, val []byte) error {
		delta := storage.BytesToCounter(val)
		counterMetric := model.Metric{ID: name, MType: counter, Delta: &delta}
		return postMetric(url, counterMetric, client)
	})
	if err != nil {
		fmt.Printf("Failed to send counter metric; %s\n", err)
		return
	}
	counterStorage.Set("PollCount", storage.CounterToBytes(storage.Counter(0)))
}
