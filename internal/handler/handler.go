package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

var defaultHeader = "text/plain; charset=utf-8"

// gauge handlers

func GaugePostHandler(gaugeStorage *storage.MemStorage[string, float64]) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", defaultHeader)
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		valueFloat, err := strconv.ParseFloat(chi.URLParam(request, "value"), 64)
		if err != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}
		gaugeStorage.Set(name, valueFloat)
		responseWriter.WriteHeader(http.StatusOK)
	}
}

// counter handlers

func CounterPostHandler(counterStorage *storage.MemStorage[string, int64]) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", defaultHeader)
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		valueInt, err := strconv.ParseInt(chi.URLParam(request, "value"), 10, 64)
		if err != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		currentValue, _ := counterStorage.Get(name)
		valueInt = currentValue + valueInt
		counterStorage.Set(name, valueInt)
		responseWriter.WriteHeader(http.StatusOK)
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
