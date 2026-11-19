package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	defaultURL             = "localhost:8080"
	defaultPolling         = 2
	defaultReport          = 10
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics.json"
)

type Config struct {
	URL     string
	Report  int
	Polling int
}

type ServerConfig struct {
	URL             string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

func NewConfig(args []string) *Config {
	urlFlag, reportFlag, pollingFlag := parseFlags(args)
	url, report, polling := setVars(urlFlag, reportFlag, pollingFlag)

	return &Config{
		URL:     url,
		Report:  report,
		Polling: polling,
	}
}

func NewServerConfig(args []string) *ServerConfig {
	urlFlag, intervalFlag, filePathFlag, restoreFlag := parseServerFlags(args)
	url, interval, filePath, restore := setServerVars(urlFlag, intervalFlag, filePathFlag, restoreFlag)

	return &ServerConfig{
		URL:             url,
		StoreInterval:   interval,
		FileStoragePath: filePath,
		Restore:         restore,
	}
}

func parseFlags(args []string) (url string, report int, polling int) {
	flags := flag.NewFlagSet("agent", flag.ExitOnError)
	urlFlag := flags.String("a", defaultURL, "Server address and port")
	reportFlag := flags.Int("r", defaultReport, "Report interval")
	pollingFlag := flags.Int("p", defaultPolling, "Polling interval")
	flags.Parse(args)
	return *urlFlag, *reportFlag, *pollingFlag
}

func parseServerFlags(args []string) (url string, interval int, filePath string, restore bool) {
	flags := flag.NewFlagSet("server", flag.ExitOnError)
	urlFlag := flags.String("a", defaultURL, "Server address and port")
	intervalFlag := flags.Int("i", defaultStoreInterval, "Store interval")
	filePathFlag := flags.String("f", defaultFileStoragePath, "File storage path")
	restoreFlag := flags.Bool("r", true, "Restore metrics from file")
	flags.Parse(args)
	return *urlFlag, *intervalFlag, *filePathFlag, *restoreFlag
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

func setServerVars(url string, interval int, filePath string, restore bool) (string, int, string, bool) {
	address := url
	storeInterval := interval
	fileStoragePath := filePath
	restoreMetrics := restore

	if addressENV := os.Getenv("ADDRESS"); addressENV != "" {
		address = addressENV
	}
	if intervalENV := os.Getenv("STORE_INTERVAL"); intervalENV != "" {
		if intervalIntENV, err := strconv.Atoi(intervalENV); err == nil {
			storeInterval = intervalIntENV
		} else {
			fmt.Printf("Incorrect STORE_INTERVAL=%q from ENV, instead using flag value %d\n", intervalENV, storeInterval)
		}
	}
	if filePathENV := os.Getenv("FILE_STORAGE_PATH"); filePathENV != "" {
		fileStoragePath = filePathENV
	}
	if restoreENV := os.Getenv("RESTORE"); restoreENV != "" {
		if restoreBoolENV, err := strconv.ParseBool(restoreENV); err == nil {
			restoreMetrics = restoreBoolENV
		} else {
			fmt.Printf("Incorrect RESTORE=%q from ENV, instead using flag value %v\n", restoreENV, restoreMetrics)
		}
	}

	if storeInterval < 0 {
		storeInterval = defaultStoreInterval
	}
	return address, storeInterval, fileStoragePath, restoreMetrics
}
