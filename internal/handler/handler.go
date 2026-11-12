package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

var (
	defaultHeader = "text/plain; charset=utf-8"
	htmlHeader    = "text/html; charset=utf-8"
)

type Repository interface {
	Get(name string) ([]byte, bool)
	Set(name string, val []byte)
	Keys() []string
	Iter(ctx context.Context, fn func(name string, val []byte) error) error
	Size() int
	Delete(name string)
}

// all metrics on "ADDRESS:PORT/"

func GetAllMetricsHandler(gaugeStorage Repository, counterStorage Repository) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", htmlHeader)
		page := GetMainPage(gaugeStorage, counterStorage)
		_, _ = io.WriteString(responseWriter, page)
		responseWriter.WriteHeader(http.StatusOK)
	}
}

// gauge handlers

func GaugePostHandler(gaugeStorage Repository) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", defaultHeader)
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		valueFloat, err := storage.ParseGauge(chi.URLParam(request, "value"))
		if err != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}
		gaugeStorage.Set(name, storage.GaugeToBytes(valueFloat))
		responseWriter.WriteHeader(http.StatusOK)
	}
}

func GaugeGetHandler(gaugeStorage Repository) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		if gaugeValue, ok := gaugeStorage.Get(name); ok {
			gaugeValueStr := storage.BytesToGauge(gaugeValue).String()
			responseWriter.Header().Set("Content-Type", defaultHeader)
			responseWriter.WriteHeader(http.StatusOK)
			_, err := io.WriteString(responseWriter, gaugeValueStr)
			if err != nil {
				return
			}
		}
		responseWriter.WriteHeader(http.StatusNotFound)
	}
}

// counter handlers

func CounterPostHandler(counterStorage Repository) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", defaultHeader)
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		valueInt, err := storage.ParseCounter(chi.URLParam(request, "value"))
		if err != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		if prevValue, ok := counterStorage.Get(name); ok {
			valueInt += storage.BytesToCounter(prevValue)
		}
		counterStorage.Set(name, storage.CounterToBytes(valueInt))
		responseWriter.WriteHeader(http.StatusOK)
	}
}

func CounterGetHandler(counterStorage Repository) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		if counterValue, ok := counterStorage.Get(name); ok {
			counterValueStr := storage.BytesToCounter(counterValue).String()
			responseWriter.Header().Set("Content-Type", defaultHeader)
			responseWriter.WriteHeader(http.StatusOK)
			_, err := io.WriteString(responseWriter, counterValueStr)
			if err != nil {
				return
			}
		}
		responseWriter.WriteHeader(http.StatusNotFound)
	}
}

// another handlers (e.g. for fails)

func FailPostHandler() func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", defaultHeader)
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
	}
}

// support functions

func getName(request *http.Request) string {
	return chi.URLParam(request, "name")
}
