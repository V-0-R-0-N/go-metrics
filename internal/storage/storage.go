package storage

import (
	"encoding/json"
	"fmt"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

type (
	gauge   float64
	counter int64
)
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

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

func CollectData(data Storage, PollCount *int, Mutex *sync.Mutex) error {

	res := runtime.MemStats{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	runtime.ReadMemStats(&res)

	Mutex.Lock()

	randomV := Float64ToGauge(float64(r.Intn(1000))) *
		Float64ToGauge(rand.Float64())
	data.PutGauge("RandomValue", randomV)
	*PollCount++

	data.PutGauge("Alloc", Float64ToGauge(float64(res.Alloc)))
	data.PutGauge("BuckHashSys", Float64ToGauge(float64(res.BuckHashSys)))
	data.PutGauge("Frees", Float64ToGauge(float64(res.Frees)))
	data.PutGauge("GCCPUFraction", Float64ToGauge(res.GCCPUFraction)) //float64
	data.PutGauge("GCSys", Float64ToGauge(float64(res.GCSys)))
	data.PutGauge("HeapAlloc", Float64ToGauge(float64(res.HeapAlloc)))
	data.PutGauge("HeapIdle", Float64ToGauge(float64(res.HeapIdle)))
	data.PutGauge("HeapInuse", Float64ToGauge(float64(res.HeapInuse)))
	data.PutGauge("HeapObjects", Float64ToGauge(float64(res.HeapObjects)))
	data.PutGauge("HeapReleased", Float64ToGauge(float64(res.HeapReleased)))
	data.PutGauge("HeapSys", Float64ToGauge(float64(res.HeapSys)))
	data.PutGauge("LastGC", Float64ToGauge(float64(res.LastGC)))
	data.PutGauge("Lookups", Float64ToGauge(float64(res.Lookups)))
	data.PutGauge("MCacheInuse", Float64ToGauge(float64(res.MCacheInuse)))
	data.PutGauge("MCacheSys", Float64ToGauge(float64(res.MCacheSys)))
	data.PutGauge("MSpanInuse", Float64ToGauge(float64(res.MSpanInuse)))
	data.PutGauge("MSpanSys", Float64ToGauge(float64(res.MSpanSys)))
	data.PutGauge("Mallocs", Float64ToGauge(float64(res.Mallocs)))
	data.PutGauge("NextGC", Float64ToGauge(float64(res.NextGC)))
	data.PutGauge("NumForcedGC", Float64ToGauge(float64(res.NumForcedGC)))
	data.PutGauge("NumGC", Float64ToGauge(float64(res.NumGC)))
	data.PutGauge("OtherSys", Float64ToGauge(float64(res.OtherSys)))
	data.PutGauge("PauseTotalNs", Float64ToGauge(float64(res.PauseTotalNs)))
	data.PutGauge("StackInuse", Float64ToGauge(float64(res.StackInuse)))
	data.PutGauge("StackSys", Float64ToGauge(float64(res.StackSys)))
	data.PutGauge("Sys", Float64ToGauge(float64(res.Sys)))
	data.PutGauge("TotalAlloc", Float64ToGauge(float64(res.TotalAlloc)))

	Mutex.Unlock()
	// В коде будет логика обработки и возврата ошибок
	return nil
}

func sendGauge(data Storage, addr *flags.NetAddress, name string) {

	value := float64(data.GetGauge(name))
	metrics := Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
	body, err := json.Marshal(metrics)
	if err != nil {
		log.Fatalln(err)
	}
	//url := fmt.Sprintf("http://%s/update/gauge/%s/%v", addr.String(), name, value)
	url := fmt.Sprintf("http://%s/update/", addr.String())
	r := strings.NewReader(string(body))
	//fmt.Println(url) // Для тестов
	//resp, err := http.Post(url, "text/plain", r)
	resp, err := http.Post(url, "application/json", r)
	if err != nil || resp.StatusCode != http.StatusOK {
		//fmt.Println("Bad response", name, value) // Для теста
		return
	}
	defer resp.Body.Close()
	//body, _ = io.ReadAll(resp.Body) // для теста
	//fmt.Println(string(body))

}

func sendCounter(addr *flags.NetAddress, PollCount *int) {

	value := int64(*PollCount)
	metrics := Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &value,
	}
	body, err := json.Marshal(metrics)
	if err != nil {
		log.Fatalln(err)
	}
	r := strings.NewReader(string(body))
	url := fmt.Sprintf("http://%s/update/", addr.String())
	//url := fmt.Sprintf("http://%s/update/counter/PollCount/%v", addr.String(), IntToCounter(*PollCount))
	resp, err := http.Post(url, "application/json", r)
	//resp, err := http.Post(url, "text/plain", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		//fmt.Println("Bad response", "PollCount", PollCount) // Для теста
		return
	}
	defer resp.Body.Close()
	//body, _ = io.ReadAll(resp.Body) // для теста
	//fmt.Println(string(body))
}

func SendData(data Storage, addr *flags.NetAddress, PollCount *int, Mutex *sync.Mutex) error {

	Mutex.Lock()

	for name := range data.GetStorage().Gauge {
		sendGauge(data.GetStorage(), addr, name)
	}
	sendCounter(addr, PollCount)
	Mutex.Unlock()
	// В коде будет логика обработки и возврата ошибок
	return nil
}
