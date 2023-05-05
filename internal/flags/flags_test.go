package flags

import (
	"log"
	"testing"
	"time"
)

func TestNetAddressSet(t *testing.T) {
	addr := NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	type hp struct {
		host string
		port int
	}
	tests := []struct {
		name  string
		data  string
		want  hp
		error bool
	}{
		{
			name: "Test 1",
			data: "localhost:8081",
			want: hp{
				host: "localhost",
				port: 8081,
			},
			error: false,
		},
		{
			name:  "Test 2 (wrong port)",
			data:  "localhost:80810",
			want:  hp{},
			error: true,
		},
		{
			name: "Test 3 (just port)",
			data: ":8082",
			want: hp{
				host: "localhost",
				port: 8082,
			},
			error: false,
		},
		{
			name:  "Test 4 (negative value)",
			data:  ":-8082",
			want:  hp{},
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := addr.Set(test.data)
			if err != nil && !test.error {
				log.Fatalf("Error in test: \"%s [%s]\"\n", test.name, err)
			}
			if (addr.Host != test.want.host || addr.Port != test.want.port) && !test.error {
				log.Fatalf("Error in test: \"%s [%s]\"\n", test.name, err)
			}
		})
	}
}

func TestPollIntervalSet(t *testing.T) {
	poll := Poll{
		Interval: 2,
	}

	tests := []struct {
		name  string
		data  string
		want  time.Duration
		error bool
	}{
		{
			name:  "Test 1",
			data:  "3",
			want:  3 * time.Second,
			error: false,
		},
		{
			name:  "Test 2 (0 value)",
			data:  "0",
			error: true,
		},
		{
			name:  "Test 3 (negative value)",
			data:  "-11",
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := poll.Set(test.data)
			if err != nil && !test.error {
				log.Fatalf("Error in test: \"%s [%s]\"\n", test.name, err)
			}
			if poll.Interval != test.want && !test.error {
				log.Fatalf("Error in test: \"%s [%s]\"\n", test.name, err)
			}
		})
	}
}

func TestReportSet(t *testing.T) {
	report := Report{
		Interval: 10,
	}

	tests := []struct {
		name  string
		data  string
		want  time.Duration
		error bool
	}{
		{
			name:  "Test 1",
			data:  "3",
			want:  3 * time.Second,
			error: false,
		},
		{
			name:  "Test 2 (0 value)",
			data:  "0",
			error: true,
		},
		{
			name:  "Test 3 (negative value)",
			data:  "-11",
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := report.Set(test.data)
			if err != nil && !test.error {
				log.Fatalf("Error in test: \"%s [%s]\"\n", test.name, err)
			}
			if report.Interval != test.want && !test.error {
				log.Fatalf("Error in test: \"%s [%s]\"\n", test.name, err)
			}
		})
	}
}
