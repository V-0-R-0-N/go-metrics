package storage

import (
	"encoding/json"
	"fmt"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"github.com/V-0-R-0-N/go-metrics.git/internal/middlware/compress"
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
	FileR   *flags.FileR
}

func NewStorage() Storage {
	return &memStorage{
		Gauge:   make(map[string]gauge),
		Counter: make(map[string]counter),
	}
}

func (m *memStorage) PutGauge(name string, value gauge) {
	m.Gauge[name] = value
	if m.FileR != nil && m.FileR.Synchro {
		n := float64(value)
		err := SaveData(compress.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &n,
		}, &m.FileR.File)
		if err != nil {
			//TODO
			log.Fatal(err)
		}
	}
}

func (m *memStorage) PutCounter(name string, value counter) {
	m.Counter[name] += value
	if m.FileR != nil && m.FileR.Synchro {
		n := int64(value)
		err := SaveData(compress.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &n,
		}, &m.FileR.File)
		if err != nil {
			//TODO
			log.Fatal(err)
		}
	}
}

// TODO обсудить с ментором как избежать цикличности в этом случае

func SaveData(metrics compress.Metrics, f *flags.OsFile) error {
	byteArr, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	byteArr = append(byteArr, byte('\n'))
	err = writeDataToFile(byteArr, f)
	if err != nil {
		return err
	}
	return nil
}

func writeDataToFile(data []byte, f *flags.OsFile) error {
	_, err := f.File.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// TODO end

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

func sendGauge(client *http.Client, data Storage, addr *flags.NetAddress, name string) error {

	value := float64(data.GetGauge(name))
	metrics := compress.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}

	url := fmt.Sprintf("http://%s/update/", addr.String())

	if err := compress.SendPostGzipJSON(client, &metrics, &url); err != nil {
		return err
	}
	return nil
}

func sendCounter(client *http.Client, addr *flags.NetAddress, PollCount *int) error {

	value := int64(*PollCount)
	metrics := compress.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &value,
	}

	url := fmt.Sprintf("http://%s/update/", addr.String())

	if err := compress.SendPostGzipJSON(client, &metrics, &url); err != nil {
		return err
	}
	return nil
}

func SendData(client *http.Client, data Storage, addr *flags.NetAddress, PollCount *int, Mutex *sync.Mutex) error {

	Mutex.Lock()

	for name := range data.GetStorage().Gauge {
		if err := sendGauge(client, data, addr, name); err != nil {
			return err
		}
	}
	if err := sendCounter(client, addr, PollCount); err != nil {
		return err
	}
	Mutex.Unlock()

	return nil
}
