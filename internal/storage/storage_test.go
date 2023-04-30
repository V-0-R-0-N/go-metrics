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
	mem := New()
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
	mem := New()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem.PutCounter(test.data.name, IntToCounter(test.data.value))
			if mem.GetCounter(test.data.name) != IntToCounter(test.want) {
				log.Fatalf("Error with name %s and value %v\n", test.data.name, IntToCounter(test.want))
			}
		})
	}
}

//func TestPut(t *testing.T) {
//	type data struct {
//		name   string
//		value1 float64
//		value2 int
//	}
//	tests := []struct {
//		name string
//		data
//		want1 float64
//		want2 int
//	}{
//		{
//			name: "Simple test 1 (Alloc, 1.2, 1)",
//			data: data{
//				name:   "Alloc",
//				value1: 1.2,
//				value2: 1,
//			},
//		},
//		{
//			name: "Simple test 2 (Malloc, 1.200000009, 30)",
//			data: data{
//				name:   "Malloc",
//				value1: 1.200000009,
//				value2: 30,
//			},
//		},
//		{
//			name: "Simple test 3 (Call, 100009.20000009, 999999999999)",
//			data: data{
//				name:   "Call",
//				value1: 100009.20000009,
//				value2: 999999999999,
//			},
//		},
//	}
//	mem := Mem{
//		Gauge:   map[string]gauge{},
//		Counter: map[string]counter{},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			Put(mem, test.data.name, Float64ToGauge(test.data.value1))
//			Put(mem, test.data.name, IntToCounter(test.data.value2))
//			if mem.Gauge[test.data.name] != Float64ToGauge(test.data.value1) {
//				log.Fatalf("Error Gauge with name: %s value:%v\n", test.data.name, mem.Gauge[test.data.name])
//			}
//			if mem.Counter[test.data.name] != IntToCounter(test.data.value2) {
//				log.Fatalf("Error Counter with name: %s value:%v\n", test.data.name, test.data.value2)
//			}
//		})
//	}
//}
