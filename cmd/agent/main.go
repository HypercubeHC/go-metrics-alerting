package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/agent"
)

var (
	DefaultPollInterval   = 2 * time.Second
	DefaultReportInterval = 10 * time.Second
)

func main() {
	url := flag.String("a", "localhost:8080", "Server address and port")
	interval := flag.Int("r", 10, "Report interval")
	polling := flag.Int("p", 2, "Polling interval")
	flag.Parse()
	fmt.Printf("Server: http://%s, Polling interval: %d sec, Reportinterval: %d sec.\n",
		*url, *polling, *interval)

	var poll time.Duration
	if *polling == 0 {
		poll = DefaultPollInterval
	} else {
		poll = time.Duration(*polling) * time.Second
	}

	var report time.Duration
	if *interval == 0 {
		report = DefaultReportInterval
	} else {
		report = time.Duration(*interval) * time.Second
	}
	agent.Run(poll, report)
}
