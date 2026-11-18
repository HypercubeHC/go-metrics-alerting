package server

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/HypercubeHC/go-metrics-alerting.git/internal/model"
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type FileSaver struct {
	mx             sync.Mutex
	filePath       string
	gaugeStorage   *storage.MemStorage[string, []byte]
	counterStorage *storage.MemStorage[string, []byte]
}

func NewFileSaver(filePath string, gaugeStrg, counterStrg *storage.MemStorage[string, []byte]) *FileSaver {
	return &FileSaver{filePath: filePath, gaugeStorage: gaugeStrg, counterStorage: counterStrg}
}

func SaveMetrics(saver *FileSaver, sugar *zap.SugaredLogger, when string) {
	if saver == nil {
		return
	}
	if err := saver.Save(); err != nil {
		sugar.Errorf("Failed to save metrics %s: %v", when, err)
	}
}

func RunStoreTicker(saver *FileSaver, interval int, sugar *zap.SugaredLogger) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		SaveMetrics(saver, sugar, "")
	}
}

func (fs *FileSaver) Load() error {
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	var metrics []model.Metric
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}
	for _, metric := range metrics {
		switch metric.MType {
		case GaugeType:
			if metric.Value != nil {
				fs.gaugeStorage.Set(metric.ID, storage.GaugeToBytes(storage.Gauge(*metric.Value)))
			}
		case CounterType:
			if metric.Delta != nil {
				fs.counterStorage.Set(metric.ID, storage.CounterToBytes(storage.Counter(*metric.Delta)))
			}
		}
	}
	return nil
}

func (fs *FileSaver) Save() error {
	fs.mx.Lock()
	defer fs.mx.Unlock()

	metrics := make([]model.Metric, 0, fs.gaugeStorage.Size()+fs.counterStorage.Size())
	for _, name := range fs.gaugeStorage.Keys() {
		gaugeBytes, _ := fs.gaugeStorage.Get(name)
		gaugeValue := float64(storage.BytesToGauge(gaugeBytes))
		metrics = append(metrics, model.Metric{ID: name, MType: GaugeType, Value: &gaugeValue})
	}
	for _, name := range fs.counterStorage.Keys() {
		counterBytes, _ := fs.counterStorage.Get(name)
		counterValue := int64(storage.BytesToCounter(counterBytes))
		metrics = append(metrics, model.Metric{ID: name, MType: CounterType, Delta: &counterValue})
	}

	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := fs.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpFile, fs.filePath)
}
