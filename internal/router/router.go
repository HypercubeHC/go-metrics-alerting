package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/handler"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	gaugeStorage := storage.New[string, float64]()
	counterStorage := storage.New[string, int64]()

	router.Route("/", func(router chi.Router) {
		router.Route("/update", func(router chi.Router) {
			router.Post("/gauge/{name}/{value}", handler.GaugePostHandler(gaugeStorage))
			router.Post("/counter/{name}/{value}", handler.CounterPostHandler(counterStorage))
			router.Post("/{another_type}/{name}/{value}", handler.FailPostHandler())
		})
		router.Post("/", handler.FailPostHandler())
	})

	return router
}
