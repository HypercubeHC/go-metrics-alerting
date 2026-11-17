package main

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
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/middleware"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/router"
)

func main() {
	conf := config.NewConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	server := &http.Server{
		Addr:    conf.URL,
		Handler: middleware.LogMiddleware(logger.Sugar())(middleware.CompressMiddleware(router.NewRouter())),
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Printf("Server started: http://%v.\nFor stop: press Ctrl + C", conf.URL)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			os.Exit(1)
		}
	}()

	sig := <-stop
	fmt.Printf("\nGot signal: (%v). Stopping...\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		server.Close()
	}
}
