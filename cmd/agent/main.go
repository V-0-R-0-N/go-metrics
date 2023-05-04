package main

import (
	"flag"
	"fmt"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	st "github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

var (
	//pollInterval   = 2 * time.Second
	//reportInterval = 10 * time.Second

	wg        = sync.WaitGroup{}
	PollCount = st.IntToCounter(0)
	Mutex     = sync.Mutex{}
)

var poll = flags.Poll{
	Interval: 2,
}

var report = flags.Report{
	Interval: 10,
}

var addr = flags.NetAddress{
	Host: "localhost",
	Port: 8080,
}

func init() {

	flags.Agent(&addr, &poll, &report)
}
func collectData(data st.Storage) {

	res := runtime.MemStats{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(poll.Interval)
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

func sendGauge(name string, data st.Storage) {

	value := data.GetStorage().GetGauge(name)
	url := name + "/" + fmt.Sprintf("%v", value)
	//fmt.Println(Host + "update/gauge/" + url) // Для тестов
	Host := "http://" + fmt.Sprintf("%s", addr) + "/"
	resp, err := http.Post(Host+"update/gauge/"+url, "text/plain", nil)
	if err != nil || resp.Status != "200 OK" {
		//fmt.Println("Bad response", name, value) // Для теста
		return
	}
	err = resp.Body.Close()
	if err != nil {
		return
	}
}

func sendCounter() {

	url := "update/counter/PollCount/" + fmt.Sprintf("%v", PollCount)
	Host := "http://" + fmt.Sprintf("%s", addr) + "/"
	resp, err := http.Post(Host+url, "text/plain", nil)
	if err != nil || resp.Status != "200 OK" {
		//fmt.Println("Bad response", "PollCount", PollCount) // Для теста
		return
	}
	err = resp.Body.Close()
	if err != nil {
		return
	}

}

func sendData(data st.Storage) {

	time.Sleep(report.Interval)

	Mutex.Lock()

	for name := range data.GetStorage().Gauge {
		sendGauge(name, data.GetStorage())
	}
	sendCounter()
	Mutex.Unlock()

}

func main() {

	flag.Parse()

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
