package agent

import (
	"context"
	"fmt"
	"time"

	"resty.dev/v3"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/config"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

func createStorages() map[string]*storage.MemStorage[string, []byte] {
	return map[string]*storage.MemStorage[string, []byte]{
		gauge:   storage.New[string, []byte](),
		counter: storage.New[string, []byte](),
	}
}

func Run() {
	conf := config.NewConfig()
	fmt.Printf("Agent started.\nServer: %s.\nPolling interval: %d sec.\nReport interval: %d sec.\n",
		conf.URL, conf.Polling, conf.Report)

	client := resty.New().SetTimeout(10 * time.Second)
	ctx := context.Background()
	storages := createStorages()

	go getMetrics(ctx, time.Duration(conf.Polling)*time.Second, storages)
	go sendMetrics(ctx, conf.URL, time.Duration(conf.Report)*time.Second, client, storages)

	<-ctx.Done()
}
