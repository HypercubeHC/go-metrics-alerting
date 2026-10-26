package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/router"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router.NewRouter(),
	}
	stop := make(chan os.Signal, 1)

	go func() {
		fmt.Println("Server started: http://localhost:8080.\nFor stop: press Ctrl + C")
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
