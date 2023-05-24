package main

import (
	"context"
	"flag"
	"net/http"
	"sync"
	"time"

	"github.com/V-0-R-0-N/go-metrics.git/internal/environ"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	st "github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

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
				err := st.CollectData(data, &PollCount, &Mutex)
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
		client := http.Client{}
		for {
			time.Sleep(report.Interval)
			select {
			case <-ctx.Done():
				return
			default:
				err := st.SendData(&client, data, &addr, &PollCount, &Mutex)
				if err == nil {
					counter = 0
				} else {
					//fmt.Println(err) // для тестов
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
