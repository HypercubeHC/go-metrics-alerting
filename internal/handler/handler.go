package handler

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/model"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

const (
	CounterType = "counter"
	GaugeType   = "gauge"
)

var (
	defaultHeader = "text/plain; charset=utf-8"
	htmlHeader    = "text/html; charset=utf-8"
	jsonHeader    = "application/json"
)

type Repository interface {
	Get(name string) ([]byte, bool)
	Set(name string, val []byte)
	Keys() []string
	Iter(ctx context.Context, fn func(name string, val []byte) error) error
	Size() int
	Delete(name string)
}

func getName(request *http.Request) string {
	return chi.URLParam(request, "name")
}

func checkRequestBody(request *http.Request, responseWriter http.ResponseWriter, metric *model.Metric) bool {
	if request.Body == nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return true
	}
	if err := json.NewDecoder(request.Body).Decode(&metric); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return true
	}
	return false
}

// all metrics on "ADDRESS:PORT/"

func GetAllMetricsHandler(gaugeStorage Repository, counterStorage Repository) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", htmlHeader)
		responseWriter.WriteHeader(http.StatusOK)
		page := GetMainPage(gaugeStorage, counterStorage)
		_, _ = io.WriteString(responseWriter, page)
	}
}

func ValuePostHandler(gaugeStorage Repository, counterStorage Repository) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		var metric model.Metric
		responseWriter.Header().Set("Content-Type", jsonHeader)

		requestBodyOK := checkRequestBody(request, responseWriter, &metric)
		if requestBodyOK {
			return
		}

		switch metric.MType {
		case GaugeType:
			val, ok := gaugeStorage.Get(metric.ID)
			if !ok {
				responseWriter.WriteHeader(http.StatusNotFound)
				return
			}
			value := storage.Gauge(storage.BytesToGauge(val))
			metric.Value = &value
		case CounterType:
			val, ok := counterStorage.Get(metric.ID)
			if !ok {
				responseWriter.WriteHeader(http.StatusNotFound)
				return
			}
			delta := storage.Counter(storage.BytesToCounter(val))
			metric.Delta = &delta
		default:
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(responseWriter).Encode(metric)
	}
}

func UpdatePostHandler(gaugeStorage Repository, counterStorage Repository) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", jsonHeader)

		var metric model.Metric
		requestBodyOK := checkRequestBody(request, responseWriter, &metric)
		if requestBodyOK {
			return
		}

		switch metric.MType {
		case GaugeType:
			if metric.ID == "" || metric.Value == nil {
				responseWriter.WriteHeader(http.StatusBadRequest)
				return
			}
			gaugeStorage.Set(metric.ID, storage.GaugeToBytes(storage.Gauge(*metric.Value)))
		case CounterType:
			if metric.ID == "" || metric.Delta == nil {
				responseWriter.WriteHeader(http.StatusBadRequest)
				return
			}
			delta := storage.Counter(*metric.Delta)
			if prevValue, ok := counterStorage.Get(metric.ID); ok {
				prev := storage.BytesToCounter(prevValue)
				if (prev > 0 && delta > math.MaxInt64-prev) || (prev < 0 && delta < math.MinInt64-prev) {
					responseWriter.WriteHeader(http.StatusBadRequest)
					return
				}
				delta += prev
			}
			counterStorage.Set(metric.ID, storage.CounterToBytes(delta))
			metric.Delta = &delta
		default:
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(responseWriter).Encode(metric)
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
			_, _ = io.WriteString(responseWriter, gaugeValueStr)
			return
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
			prev := storage.BytesToCounter(prevValue)
			if (prev > 0 && valueInt > math.MaxInt64-prev) || (prev < 0 && valueInt < math.MinInt64-prev) {
				responseWriter.WriteHeader(http.StatusBadRequest)
				return
			}
			valueInt += prev
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
			_, _ = io.WriteString(responseWriter, counterValueStr)
			return
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
