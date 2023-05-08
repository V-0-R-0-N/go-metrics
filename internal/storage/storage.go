package storage

import (
	"fmt"
	"strings"
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

	GetStorage() *memStorage
}

type memStorage struct {
	Gauge   map[string]gauge
	Counter map[string]counter
}

func NewStorage() Storage {
	return &memStorage{
		Gauge:   make(map[string]gauge),
		Counter: make(map[string]counter),
	}
}

func (m *memStorage) PutGauge(name string, value gauge) {
	m.Gauge[name] = value
}

func (m *memStorage) PutCounter(name string, value counter) {
	m.Counter[name] += value
}

func (m *memStorage) GetGauge(name string) gauge {
	return m.Gauge[name]
}

func (m *memStorage) GetCounter(name string) counter {
	return m.Counter[name]
}

func (m *memStorage) GetStorage() *memStorage {
	return m
}

func Float64ToGauge(v float64) gauge {
	return gauge(v)
}

func IntToCounter(v int) counter {
	return counter(v)
}

func (m *memStorage) String() string {
	res := strings.Builder{}
	res.WriteString("Gauge\n")
	for k, v := range m.Gauge {
		res.WriteString(fmt.Sprintf("%21s:\t\t\t%21v\n", k, v))
		//res += k + ": " + fmt.Sprintf("%v\n", v)
	}
	res.WriteString("\nCounter\n")
	for k, v := range m.Counter {
		res.WriteString(fmt.Sprintf("%21s:\t\t\t%21v\n", k, v))
		//res += k + ": " + fmt.Sprintf("%v\n", v)
	}
	return res.String()
}
