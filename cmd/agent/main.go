package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	st "github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

var (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	//data           = st.Memory
	wg        = sync.WaitGroup{}
	PollCount = st.IntToCounter(0)
	Host      = "http://localhost:8080/"
	Mutex     = sync.Mutex{}
)

func collectData(data st.Storage) {

	res := runtime.MemStats{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(pollInterval)
	runtime.ReadMemStats(&res)

	Mutex.Lock()

	randomV := st.Float64ToGauge(float64(r.Intn(1000))) *
		st.Float64ToGauge(rand.Float64())
	data.PutGauge("RandomValue", randomV)
	PollCount++

	data.PutGauge("Alloc", st.Float64ToGauge(float64(res.Alloc)))
	data.PutGauge("BuckHashSys", st.Float64ToGauge(float64(res.BuckHashSys)))
	data.PutGauge("Frees", st.Float64ToGauge(float64(res.Frees)))
	data.PutGauge("GCCPUFraction", st.Float64ToGauge(res.GCCPUFraction)) //float64
	data.PutGauge("GCSys", st.Float64ToGauge(float64(res.GCSys)))
	data.PutGauge("HeapAlloc", st.Float64ToGauge(float64(res.HeapAlloc)))
	data.PutGauge("HeapIdle", st.Float64ToGauge(float64(res.HeapIdle)))
	data.PutGauge("HeapInuse", st.Float64ToGauge(float64(res.HeapInuse)))
	data.PutGauge("HeapObjects", st.Float64ToGauge(float64(res.HeapObjects)))
	data.PutGauge("HeapReleased", st.Float64ToGauge(float64(res.HeapReleased)))
	data.PutGauge("HeapSys", st.Float64ToGauge(float64(res.HeapSys)))
	data.PutGauge("LastGC", st.Float64ToGauge(float64(res.LastGC)))
	data.PutGauge("Lookups", st.Float64ToGauge(float64(res.Lookups)))
	data.PutGauge("MCacheInuse", st.Float64ToGauge(float64(res.MCacheInuse)))
	data.PutGauge("MCacheSys", st.Float64ToGauge(float64(res.MCacheSys)))
	data.PutGauge("MSpanInuse", st.Float64ToGauge(float64(res.MSpanInuse)))
	data.PutGauge("MSpanSys", st.Float64ToGauge(float64(res.MSpanSys)))
	data.PutGauge("Mallocs", st.Float64ToGauge(float64(res.Mallocs)))
	data.PutGauge("NextGC", st.Float64ToGauge(float64(res.NextGC)))
	data.PutGauge("NumForcedGC", st.Float64ToGauge(float64(res.NumForcedGC)))
	data.PutGauge("NumGC", st.Float64ToGauge(float64(res.NumGC)))
	data.PutGauge("OtherSys", st.Float64ToGauge(float64(res.OtherSys)))
	data.PutGauge("PauseTotalNs", st.Float64ToGauge(float64(res.PauseTotalNs)))
	data.PutGauge("StackInuse", st.Float64ToGauge(float64(res.StackInuse)))
	data.PutGauge("StackSys", st.Float64ToGauge(float64(res.StackSys)))
	data.PutGauge("Sys", st.Float64ToGauge(float64(res.Sys)))
	data.PutGauge("TotalAlloc", st.Float64ToGauge(float64(res.TotalAlloc)))

	Mutex.Unlock()
}

func sendData(data st.Storage) {

	time.Sleep(reportInterval)
	resp := http.Response{
		Body: io.NopCloser(bytes.NewBufferString("Hello World")),
	}

	Mutex.Lock()

	for name, value := range data.GetStorage().Gauge {
		url := name + "/" + fmt.Sprintf("%v", value)
		//fmt.Println(Host + "update/gauge/" + url) // Для тестов
		resp, err := http.Post(Host+"update/gauge/"+url, "text/plain", nil)
		if err != nil || resp.Status != "200 OK" {
			fmt.Println("Bad response", name, value) // Для теста
		}
	}
	for i := 0; i < 1; i++ {
		url := "update/counter/PollCount/" + fmt.Sprintf("%v", PollCount)
		resp, err := http.Post(Host+url, "text/plain", nil)
		if err != nil || resp.Status != "200 OK" {
			fmt.Println("Bad response", "PollCount", PollCount) // Для теста
		}
	}
	Mutex.Unlock()
	_ = resp.Body.Close()
}

func main() {

	wg.Add(2)
	data := st.New()
	go func() {
		defer wg.Done()
		for {
			collectData(data)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			sendData(data)
		}
	}()

	wg.Wait()
}
