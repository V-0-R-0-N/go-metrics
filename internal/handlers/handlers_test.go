package handlers

import (
	"fmt"
	st "github.com/V-0-R-0-N/go-metrics.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateValidator(t *testing.T) {
	tests := []struct {
		name   string
		splURI []string
		req    http.Request
		want   bool
	}{
		{
			name: "Simple test 1",
			splURI: []string{
				"update",
				"gauge",
				"Alloc",
				"11.7",
			},
			req: http.Request{
				Header: map[string][]string{
					"Content-Type": {"text/plain"},
				},
			},
			want: true,
		},
		{
			name: "Simple test 2 have no parameter name",
			splURI: []string{
				"update",
				"gauge",
			},
			req: http.Request{
				Header: map[string][]string{
					"Content-Type": {"text/plain"},
				},
			},
			want: false,
		},
		{
			name: "Simple test 3 wrong metric type",
			splURI: []string{
				"update",
				"wrong",
				"Alloc",
				"11.7",
			},
			req: http.Request{
				Header: map[string][]string{
					"Content-Type": {"text/plain"},
				},
			},
			want: false,
		},
		{
			name: "Simple test 3 wrong metric type",
			splURI: []string{
				"update",
				"gaige",
				"Alloc",
				"11.7",
			},
			req: http.Request{
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !updateValidator(test.splURI, &test.req) != test.want {
				log.Fatalf("Error with test: %s\n", test.name)
			}
		})
	}
}

func TestNewHandlerStorage(t *testing.T) {
	t.Run("Simple test", func(t *testing.T) {
		Hand := NewHandlerStorage(st.New())
		if Hand == nil {
			panic("aaa")
		}
	})
}

func TestUpdateMetrics(t *testing.T) {
	type want struct {
		code        int
		contentType string
		req         string
	}
	tests := []struct {
		name string
		want want
	}{
		// TODO: Add test cases.
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "text/plain",
				req:         "/update/gauge/Alloc/123.6",
			},
		},
		{
			name: "positive test #2",
			want: want{
				code:        200,
				contentType: "text/plain",
				req:         "/update/counter/Alloc/123",
			},
		},
		{
			name: "negative test #3",
			want: want{
				code:        404,
				contentType: "text/plain",
				req:         "/update/counter",
			},
		},
		{
			name: "negative test #4",
			want: want{
				code:        400,
				contentType: "text/plain",
				req:         "/update/",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.want.req, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h := NewHandlerStorage(st.New())
			h.UpdateMetrics(w, request)
			res := w.Result()
			fmt.Println(test.name, res.StatusCode)
			assert.Equal(t, res.StatusCode, test.want.code)
		})
	}
}

func TestBadRequest(t *testing.T) {
	type want struct {
		code        int
		contentType string
		req         string
	}
	tests := []struct {
		name string
		want want
	}{
		// TODO: Add test cases.
		{
			name: "positive test #1",
			want: want{
				code:        400,
				contentType: "text/plain",
				req:         "/update/",
			},
		},
		{
			name: "positive test #2",
			want: want{
				code:        400,
				contentType: "text/plain",
				req:         "/update/counter/Alloc/123/11",
			},
		},
		{
			name: "positive test #3",
			want: want{
				code:        400,
				contentType: "text/plain",
				req:         "/",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.want.req, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			BadRequest(w, request)
			res := w.Result()
			assert.Equal(t, res.StatusCode, test.want.code)
		})
	}
}
