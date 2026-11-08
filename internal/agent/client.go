package agent

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	"resty.dev/v3"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

const (
	defaultURL     = "http://localhost:8080"
	defaultPolling = 2
	defaultReport  = 10
	gauge          = "gauge"
	counter        = "counter"
)

// использовать sync.Mutex мне посоветовала нейронка, я бы сама не догадалась

type metricStorage struct {
	mx sync.Mutex
	Repository
}

func createStorages() map[string]*metricStorage {
	return map[string]*metricStorage{
		gauge: {
			mx:         sync.Mutex{},
			Repository: storage.New[string, []byte](),
		},
		counter: {
			mx:         sync.Mutex{},
			Repository: storage.New[string, []byte](),
		},
	}
}

func Run() {
	url := flag.String("a", defaultURL, "Server address and port")
	report := flag.Int("r", defaultReport, "Report interval")
	polling := flag.Int("p", defaultPolling, "Polling interval")
	flag.Parse()
	fmt.Printf("Agent started.\nServer: %s.\nPolling interval: %d sec.\nReport interval: %d sec.\n",
		*url, *polling, *report)

	client := resty.New()
	ctx := context.Background()
	storages := createStorages()

	go getMetrics(ctx, time.Duration(*polling)*time.Second, storages)
	go sendMetrics(ctx, *url, time.Duration(*report)*time.Second, client, storages)

	<-ctx.Done()
}
