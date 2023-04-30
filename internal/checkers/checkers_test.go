package checkers

import (
	"log"
	"net/http"
	"testing"
)

func TestCheckMetricType(t *testing.T) {
	tests := []struct {
		name string
		data string

		want bool
	}{
		{
			name: "Test 1 (gauge)",
			data: "gauge",
			want: true,
		},
		{
			name: "Test 2 (counter)",
			data: "counter",
			want: true,
		},
		{
			name: "Test 3 (wrong data)",
			data: "wrong data",
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if CheckMetricType(test.data) != test.want {
				log.Fatalf("Vrong ansver with data: \"%s\"\n", test.data)
			}
		})
	}
}

func TestCheckContentType(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{
			name: "Simple test 1 (True)",
			req: &http.Request{
				Header: map[string][]string{
					"Content-Type": {"text/plain"},
				},
			},
			want: true,
		},
		{
			name: "Simple test 2 (True (empty))",
			req:  &http.Request{},
			want: true,
		},
		{
			name: "Simple test 3 (False)",
			req: &http.Request{
				Header: map[string][]string{
					"Content-Type": {
						"application/json",
					},
				},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if CheckContentType(test.req) != test.want {
				log.Fatalf("Wrong Content-Type check test: %s\n", test.name)
			}
		})
	}
}
