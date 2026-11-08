package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/router"
)

const defaultURL = "http://localhost:8080"

func main() {
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	url := flag.String("a", defaultURL, "Server address and port")
	flag.Parse()
	server := &http.Server{
		Addr:    *url,
		Handler: router.NewRouter(),
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Printf("Server started: http://%v.\nFor stop: press Ctrl + C", *url)
		if err := server.ListenAndServe(); err != nil {
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
