package model

import (
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

type Metric struct {
	ID    string           `json:"id"`
	MType string           `json:"type"`
	Delta *storage.Counter `json:"delta,omitempty"`
	Value *storage.Gauge   `json:"value,omitempty"`
	Hash  string           `json:"hash,omitempty"`
}
