package agent

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"resty.dev/v3"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

func Run(poll time.Duration, report time.Duration) {
	gaugeStorage := storage.New[string, gaugeValue]()
	counterStorage := storage.New[string, counterValue]()

	pollTicker := time.NewTicker(poll)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(report)
	defer reportTicker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	client := resty.New()
	for {
		select {
		case <-pollTicker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			getGaugeMetrics(gaugeStorage)
			getCounterMetrics(counterStorage)

		case <-reportTicker.C:
			sendGauge(gaugeStorage, client)
			sendCounter(counterStorage, client)

		case sig := <-stop:
			fmt.Printf("Agent got signal: %v.\n", sig)
			return
		}
	}
}
