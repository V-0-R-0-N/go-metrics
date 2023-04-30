package storage

import (
	"fmt"
)

type (
	gauge   float64
	counter int64
)

type Storage interface {
	PutGauge(name string, value gauge)
	PutCounter(name string, value counter)

	GetGauge(name string) gauge
	GetCounter(name string) counter

	GetStorage() MemStorage
}

type MemStorage struct {
	Gauge   map[string]gauge
	Counter map[string]counter
}

func New() Storage {
	return MemStorage{
		Gauge:   make(map[string]gauge),
		Counter: make(map[string]counter),
	}
}

func (m MemStorage) PutGauge(name string, value gauge) {
	m.Gauge[name] = value
}

func (m MemStorage) PutCounter(name string, value counter) {
	m.Counter[name] += value
}

func (m MemStorage) GetGauge(name string) gauge {
	return m.Gauge[name]
}
func (m MemStorage) GetCounter(name string) counter {
	return m.Counter[name]
}

func (m MemStorage) GetStorage() MemStorage {
	return m
}

func Float64ToGauge(v float64) gauge {
	return gauge(v)
}

func IntToCounter(v int) counter {
	return counter(v)
}

func (m MemStorage) String() string {
	res := ""
	res += "Gauge\n"
	for k, v := range m.Gauge {
		res += k + ": " + fmt.Sprintf("%v\n", v)
	}
	res += "\nCounter\n"
	for k, v := range m.Counter {
		res += k + ": " + fmt.Sprintf("%v\n", v)
	}
	return res
}
