package filer

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"github.com/V-0-R-0-N/go-metrics.git/internal/middlware/compress"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
	"log"
	"os"
	"time"
)

func NewFile(filename string) *os.File {
	fl, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return fl
}

func FilerInit(FileR *flags.FileRestore) {
	if FileR.Path.Data != "" {
		FileR.Restore = true
		FileR.File = NewFile(FileR.Path.Data)
		//defer Close(&FileR.File)
	}
	if FileR.Interval.Data == 0 {
		FileR.Synchro = true
	}
}

func StartRestore(ctx context.Context, st storage.Storage, FileR *flags.FileRestore) {
	if FileR.FileRestore.Data {
		err := RestoreData(st, FileR.File)
		if err != nil {
			log.Fatal(err)
		}
	}
	if !FileR.Synchro {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(FileR.Interval.Data)
					SaveAllData(st, FileR.File)
				}
			}
		}()
	} else {
		st.GetStorage().FileR = FileR
	}
}

func SaveAllData(data storage.Storage, f *os.File) error {
	allData := data.GetStorage()
	metrics := compress.Metrics{}
	metrics.MType = "gauge"
	for k, v := range allData.Gauge {
		metrics.ID = k
		n := float64(v)
		metrics.Value = &n

		byteArr, err := json.Marshal(metrics)
		if err != nil {
			return err
		}
		byteArr = append(byteArr, byte('\n'))
		err = writeDataToFile(byteArr, f)
		if err != nil {
			return err
		}
	}
	metrics = compress.Metrics{}
	metrics.MType = "counter"
	for k, v := range allData.Counter {
		metrics.ID = k
		n := int64(v)
		metrics.Delta = &n

		byteArr, err := json.Marshal(metrics)
		if err != nil {
			return err
		}
		byteArr = append(byteArr, byte('\n'))
		err = writeDataToFile(byteArr, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveData(metrics compress.Metrics, f *os.File) error {
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

func writeDataToFile(data []byte, f *os.File) error {
	_, err := f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func RestoreData(data storage.Storage, f *os.File) error {

	metrics := compress.Metrics{}
	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		//fmt.Println(scanner.Text()) // для теста
		err := json.Unmarshal(scanner.Bytes(), &metrics)
		if err != nil {
			return err
		}
		if metrics.MType == "gauge" {
			data.PutGauge(metrics.ID, storage.Float64ToGauge(*metrics.Value))
		} else if metrics.MType == "counter" {
			data.GetStorage().Counter[metrics.ID] = storage.IntToCounter(int(*metrics.Delta))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
