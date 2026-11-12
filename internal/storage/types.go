package storage

import (
	"encoding/binary"
	"math"
	"strconv"
)

type (
	Gauge float64
)

func (gauge Gauge) String() string {
	return strconv.FormatFloat(float64(gauge), 'f', -1, 64)
}

func ParseGauge(gaugeString string) (Gauge, error) {
	gauge, err := strconv.ParseFloat(gaugeString, 64)
	return Gauge(gauge), err
}

func GaugeToBytes(gauge Gauge) []byte {
	var gaugeBytes [8]byte
	binary.LittleEndian.PutUint64(gaugeBytes[:], math.Float64bits(float64(gauge)))
	return gaugeBytes[:]
}

func BytesToGauge(gaugeBytes []byte) Gauge {
	gauge := binary.LittleEndian.Uint64(gaugeBytes)
	return Gauge(math.Float64frombits(gauge))
}

type (
	Counter int64
)

func (counter Counter) String() string {
	return strconv.FormatInt(int64(counter), 10)
}

func ParseCounter(counterString string) (Counter, error) {
	counter, err := strconv.ParseInt(counterString, 10, 64)
	return Counter(counter), err
}

func CounterToBytes(counter Counter) []byte {
	var counterBytes [8]byte
	binary.LittleEndian.PutUint64(counterBytes[:], uint64(counter))
	return counterBytes[:]
}

func BytesToCounter(counterBytes []byte) Counter {
	counter := binary.LittleEndian.Uint64(counterBytes)
	return Counter(counter)
}
