package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/handler"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	gaugeStorage := storage.New[string, []byte]()
	counterStorage := storage.New[string, []byte]()

	router.Route("/", func(router chi.Router) {
		router.Route("/update", func(router chi.Router) {
			router.Post("/gauge/{name}/{value}", handler.GaugePostHandler(gaugeStorage))
			router.Post("/counter/{name}/{value}", handler.CounterPostHandler(counterStorage))
			router.Post("/{another_type}/{name}/{value}", handler.FailPostHandler())
		})
		router.Route("/value", func(router chi.Router) {
			router.Get("/gauge/{name}", handler.GaugeGetHandler(gaugeStorage))
			router.Get("/counter/{name}", handler.CounterGetHandler(counterStorage))
		})
		router.Post("/", handler.FailPostHandler())
		router.Get("/", handler.GetAllMetricsHandler(gaugeStorage, counterStorage))
	})

	return router
}
