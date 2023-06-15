package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	mock2 "github.com/V-0-R-0-N/go-metrics.git/internal/middlware/logger/mock"
)

func handler(w http.ResponseWriter, r *http.Request) {
	response := []byte("Hello, world!")
	_, _ = w.Write(response)
}

func TestWithLogging(t *testing.T) {
	type args struct {
		h http.HandlerFunc
	}

	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "Simple test 1",
			args: args{
				h: handler,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := mock2.LoggerMock{}

			middlware := Middlware{Log: &mock}

			middlware.WithLogging(tt.args.h)

			fmt.Println("|", mock.Get(), "|")
			assert.Equal(t, tt.expected, mock.Get())
		})
	}
}
