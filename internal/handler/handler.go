package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/template"
)

var (
	defaultHeader = "text/plain; charset=utf-8"
	htmlHeader    = "text/html; charset=utf-8"
)

// all metrics on "localhost:8080/"

func GetAllMetricsHandler(gaugeStorage *storage.MemStorage[string, float64], counterStorage *storage.MemStorage[string, int64]) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", htmlHeader)
		page := template.GetMainPage(gaugeStorage, counterStorage)
		_, _ = io.WriteString(responseWriter, page)
		responseWriter.WriteHeader(http.StatusOK)
	}
}

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

func GaugeGetHandler(gaugeStorage *storage.MemStorage[string, float64]) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		if gaugeValue, ok := gaugeStorage.Get(name); ok {
			gaugeValueStr := strconv.FormatFloat(gaugeValue, 'g', -1, 64)
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

func CounterGetHandler(counterStorage *storage.MemStorage[string, int64]) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		name := getName(request)
		if len(name) == 0 {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		if counterValue, ok := counterStorage.Get(name); ok {
			counterValueStr := strconv.FormatInt(counterValue, 10)
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
