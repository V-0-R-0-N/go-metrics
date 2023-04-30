package main

import (
	st "github.com/V-0-R-0-N/go-metrics.git/internal/storage"
	"log"
	"testing"
)

func TestCollectData(t *testing.T) {
	tests := []struct {
		name string
		test string
		want bool
	}{
		{
			name: "Simple test 1(Alloc)",
			test: "Alloc",
			want: true,
		},
		{
			name: "Simple test 2(Frees)",
			test: "Frees",
			want: true,
		},
		{
			name: "Simple test 3(Wrong data)",
			test: "Wrong MemStorage",
			want: false,
		},
	}
	data := st.New()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			collectData(data)
			if _, ok := data.GetStorage().Gauge[test.test]; ok != test.want {
				log.Fatalf("Have no element: \"%s\"\n", test.test)
			}
		})
	}
}

func TestSendData(t *testing.T) { // Заглушка потому что ничего не возвращает

	t.Run("Simple Test", func(t *testing.T) {
		data := st.New()
		sendData(data)
	})
}
