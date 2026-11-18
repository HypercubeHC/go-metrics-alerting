package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/config"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/middleware"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/router"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/server"
)

func main() {
	conf := config.NewServerConfig(os.Args[1:])
	logger, sugar := server.NewLogger()
	defer logger.Sync()

	gaugeStorage, counterStorage, saver := server.SetupStorages(conf, sugar)

	serverHTTP := &http.Server{
		Addr: conf.URL,
		Handler: middleware.LogMiddleware(sugar)(
			middleware.CompressMiddleware(
				router.NewRouter(gaugeStorage, counterStorage),
			),
		),
	}

	go server.Run(serverHTTP, conf.URL)

	sig := server.WaitShutdownSignal()
	fmt.Printf("\nGot signal: (%v). Stopping...\n", sig)

	server.Shutdown(serverHTTP)
	server.SaveMetrics(saver, sugar, "on shutdown")
}
