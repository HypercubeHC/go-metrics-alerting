package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	defaultURL     = "localhost:8080"
	defaultPolling = 2
	defaultReport  = 10
)

type Config struct {
	URL     string
	Report  int
	Polling int
}

func NewConfig() *Config {
	urlFlag, reportFlag, pollingFlag := parseFlags()
	url, report, polling := setVars(urlFlag, reportFlag, pollingFlag)

	return &Config{
		URL:     url,
		Report:  report,
		Polling: polling,
	}
}

func parseFlags() (url string, report int, polling int) {
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	urlFlag := flag.String("a", defaultURL, "Server address and port")
	reportFlag := flag.Int("r", defaultReport, "Report interval")
	pollingFlag := flag.Int("p", defaultPolling, "Polling interval")
	flag.Parse()
	return *urlFlag, *reportFlag, *pollingFlag
}

func setVars(url string, report int, polling int) (string, int, int) {
	address := url
	reportInterval := report
	pollInterval := polling

	if addressENV := os.Getenv("ADDRESS"); addressENV != "" {
		address = addressENV
	}
	if reportENV := os.Getenv("REPORT_INTERVAL"); reportENV != "" {
		if reportIntENV, err := strconv.Atoi(reportENV); err == nil {
			reportInterval = reportIntENV
		} else {
			fmt.Printf("Incorrect REPORT_INTERVAL=%q from ENV, instead using flag value %d\n", reportENV, reportInterval)
		}
	}
	if pollingENV := os.Getenv("POLL_INTERVAL"); pollingENV != "" {
		if pollingIntENV, err := strconv.Atoi(pollingENV); err == nil {
			pollInterval = pollingIntENV
		} else {
			fmt.Printf("Incorrect POLL_INTERVAL=%q from ENV, instead using flag value %d\n", pollingENV, pollInterval)
		}
	}

	if reportInterval <= 0 {
		reportInterval = defaultReport
	}
	if pollInterval <= 0 {
		pollInterval = defaultPolling
	}
	return address, reportInterval, pollInterval
}
