package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/V-0-R-0-N/go-metrics.git/internal/environ"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	st "github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

func collectData(data st.Storage, PollCount *int, Mutex *sync.Mutex) error {

	res := runtime.MemStats{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	runtime.ReadMemStats(&res)

	Mutex.Lock()

	randomV := st.Float64ToGauge(float64(r.Intn(1000))) *
		st.Float64ToGauge(rand.Float64())
	data.PutGauge("RandomValue", randomV)
	*PollCount++

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
	// В коде будет логика обработки и возврата ошибок
	return nil
}

func sendGauge(data st.Storage, addr *flags.NetAddress, name string) {

	value := data.GetGauge(name)
	url := fmt.Sprintf("http://%s/update/gauge/%s/%v", addr.String(), name, value)
	//fmt.Println(url) // Для тестов
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		//fmt.Println("Bad response", name, value) // Для теста
		return
	}
	defer resp.Body.Close()

}

func sendCounter(addr *flags.NetAddress, PollCount *int) {

	url := fmt.Sprintf("http://%s/update/counter/PollCount/%v", addr.String(), st.IntToCounter(*PollCount))
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		//fmt.Println("Bad response", "PollCount", PollCount) // Для теста
		return
	}
	defer resp.Body.Close()

}

func sendData(data st.Storage, addr *flags.NetAddress, PollCount *int, Mutex *sync.Mutex) error {

	Mutex.Lock()

	for name := range data.GetStorage().Gauge {
		sendGauge(data.GetStorage(), addr, name)
	}
	sendCounter(addr, PollCount)
	Mutex.Unlock()
	// В коде будет логика обработки и возврата ошибок
	return nil
}

func main() {

	wg := sync.WaitGroup{}

	PollCount := 0
	Mutex := sync.Mutex{}

	poll := flags.Poll{
		Interval: 2 * time.Second,
	}
	report := flags.Report{
		Interval: 10 * time.Second,
	}
	addr := flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	flags.Agent(&addr, &poll, &report)
	//fmt.Println(addr, poll.Interval, report.Interval) // Для теста
	flag.Parse()
	if err := environ.Agent(&addr, &poll, &report); err != nil {
		panic(err)
	}
	//fmt.Println(addr, poll.Interval, report.Interval) // Для теста
	data := st.NewStorage()
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		counter := 0
		for {
			time.Sleep(poll.Interval)
			select {
			case <-ctx.Done():
				return
			default:
				err := collectData(data, &PollCount, &Mutex)
				if err == nil {
					counter = 0
				} else {
					counter++
					if counter == 3 {
						cancel()
					}
				}
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		counter := 0
		for {
			time.Sleep(report.Interval)
			select {
			case <-ctx.Done():
				return
			default:
				err := sendData(data, &addr, &PollCount, &Mutex)
				if err == nil {
					counter = 0
				} else {
					counter++
					if counter == 3 {
						cancel()
					}
				}
			}
		}
	}()

	wg.Wait()
}
