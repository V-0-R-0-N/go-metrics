package filer

import (
	"fmt"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
	"os"
	"testing"
)

func TestNewFile(t *testing.T) {
	type args struct {
		filename string
		f        *os.File
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Simple test create file 1",
			args: args{
				filename: "/tmp/test.json",
				f:        nil,
			},
		},
		{
			name: "Simple test create file 2",
			args: args{
				filename: "/tmp/test2.json",
				f:        nil,
			},
		},
		{
			name: "Simple test create file 3",
			args: args{
				filename: "/tmp/test3.json",
				f:        nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.f = NewFile(tt.args.filename)
			tt.args.f.Close()
		})
	}
}

func TestSaveAllData(t *testing.T) {
	type args struct {
		data storage.Storage
		f    *os.File
	}
	file, _ := os.OpenFile("/tmp/test2.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	st := storage.NewStorage()
	st.PutGauge("Data_gauge", storage.Float64ToGauge(1.2))
	st.PutGauge("Alloc", storage.Float64ToGauge(1.209))
	st.PutCounter("Data_counter", storage.IntToCounter(3))
	st.PutCounter("Counter", storage.IntToCounter(30))
	st.PutCounter("Counter2", storage.IntToCounter(90))
	st2 := storage.NewStorage()
	st2.PutGauge("Data_gauge", storage.Float64ToGauge(1.2))
	st2.PutGauge("Alloc", storage.Float64ToGauge(1.209))
	st2.PutCounter("Data_counter", storage.IntToCounter(3))
	st2.PutCounter("Counter", storage.IntToCounter(30))
	st2.PutCounter("Counter2", storage.IntToCounter(90))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Simple test 1",
			args: args{
				data: st,
				f:    file,
			},
		},
		{
			name: "Simple test 2",
			args: args{
				data: st2,
				f:    file,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveAllData(tt.args.data, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("SaveAllData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestoreData(t *testing.T) {
	type args struct {
		data storage.Storage
		f    *os.File
	}
	file, _ := os.OpenFile("/tmp/test2.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Simple test 1",
			args: args{
				data: storage.NewStorage(),
				f:    file,
			},
		},
		{
			name: "Simple test 2",
			args: args{
				data: storage.NewStorage(),
				f:    file,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RestoreData(tt.args.data, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("RestoreData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		fmt.Println(tt.args.data)
	}
}
