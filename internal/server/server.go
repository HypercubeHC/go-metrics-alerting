package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/config"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/handler"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

func Run(serverHTTP *http.Server, url string) {
	fmt.Printf("Server started: http://%v.\nFor stop: press Ctrl + C", url)
	if err := serverHTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		os.Exit(1)
	}
}

func WaitShutdownSignal() os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return <-stop
}

func Shutdown(serverHTTP *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serverHTTP.Shutdown(ctx); err != nil {
		serverHTTP.Close()
	}
}

func NewLogger() (*zap.Logger, *zap.SugaredLogger) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger, logger.Sugar()
}

func SetupStorages(conf *config.ServerConfig, sugar *zap.SugaredLogger) (gauge, counter handler.Repository, saver *FileSaver) {
	gaugeStorage, counterStorage := storage.New[string, []byte](), storage.New[string, []byte]()
	gauge, counter = gaugeStorage, counterStorage

	if conf.FileStoragePath == "" {
		return gauge, counter, nil
	}

	saver = NewFileSaver(conf.FileStoragePath, gaugeStorage, counterStorage)
	if conf.Restore {
		if err := saver.Load(); err != nil {
			sugar.Errorf("Failed to restore metrics from %s: %v", conf.FileStoragePath, err)
		}
	}

	if conf.StoreInterval == 0 {
		syncSave := func() { SaveMetrics(saver, sugar, "") }
		gauge = NewFileStorage(gaugeStorage, syncSave)
		counter = NewFileStorage(counterStorage, syncSave)
	} else {
		go RunStoreTicker(saver, conf.StoreInterval, sugar)
	}
	return gauge, counter, saver
}
