package environ

import (
	"log"
	"testing"
	"time"

	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
)

func TestServer(t *testing.T) {
	addr := flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	}

	FileR := flags.FileR{}
	t.Run("Simple test", func(t *testing.T) {
		Server(&addr, &FileR)
		if addr.String() != "localhost:8080" {
			log.Fatal("Error Test 1")
		}
	})
}

func TestAgent(t *testing.T) {
	addr := flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	poll := flags.Poll{
		Interval: 2 * time.Second,
	}
	report := flags.Report{
		Interval: 10 * time.Second,
	}
	t.Run("Simple test 1", func(t *testing.T) {
		Agent(&addr, &poll, &report)
		if addr.String() != "localhost:8080" ||
			poll.String() != "2" || report.String() != "10" {
			log.Fatal("Error Test 1")
		}
	})
}
