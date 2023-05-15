package storage

import (
	"log"
	"testing"
)

func TestPutGauge(t *testing.T) {

	type data struct {
		name  string
		value float64
	}

	tests := []struct {
		name string
		data
		want float64
	}{
		{
			name: "Simple test 1 (Alloc: 1.2)",
			data: data{
				name:  "Alloc",
				value: 1.2,
			},
			want: 1.2,
		},
		{
			name: "Simple test 2 (Malloc: 1.20000000009)",
			data: data{
				name:  "Malloc",
				value: 1.20000000009,
			},
			want: 1.20000000009,
		},
		{
			name: "Simple test 2 (Call: 390.20000000000001)",
			data: data{
				name:  "Call",
				value: 390.20000000000001,
			},
			want: 390.20000000000001,
		},
	}
	mem := NewStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem.PutGauge(test.data.name, Float64ToGauge(test.data.value))
			if mem.GetGauge(test.data.name) != Float64ToGauge(test.want) {
				log.Fatalf("Error with name %s and value %v\n", test.data.name, Float64ToGauge(test.want))
			}
		})
	}
}

func TestPutCounter(t *testing.T) {

	type data struct {
		name  string
		value int
	}

	tests := []struct {
		name string
		data
		want int
	}{
		{
			name: "Simple test 1 (Alloc: 1)",
			data: data{
				name:  "Alloc",
				value: 1,
			},
			want: 1,
		},
		{
			name: "Simple test 2 (Malloc: 201)",
			data: data{
				name:  "Malloc",
				value: 201,
			},
			want: 201,
		},
		{
			name: "Simple test 2 (Call: 9999999999999999)",
			data: data{
				name:  "Call",
				value: 9999999999999999,
			},
			want: 9999999999999999,
		},
	}
	mem := NewStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem.PutCounter(test.data.name, IntToCounter(test.data.value))
			if mem.GetCounter(test.data.name) != IntToCounter(test.want) {
				log.Fatalf("Error with name %s and value %v\n", test.data.name, IntToCounter(test.want))
			}
		})
	}
}
