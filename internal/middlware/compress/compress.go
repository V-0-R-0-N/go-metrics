package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"syscall"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		//fmt.Println(acceptEncoding) // для тестов
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
			ow.Header().Set("Content-Encoding", "gzip")
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		//fmt.Println(contentEncoding) // для тнстов
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ow, r)
	})
}

//func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
//		// который будем передавать следующей функции
//		ow := w
//
//		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
//		acceptEncoding := r.Header.Get("Accept-Encoding")
//		supportsGzip := strings.Contains(acceptEncoding, "gzip")
//		//fmt.Println(acceptEncoding) // для тестов
//		if supportsGzip {
//			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
//			cw := newCompressWriter(w)
//			// меняем оригинальный http.ResponseWriter на новый
//			ow = cw
//			// не забываем отправить клиенту все сжатые данные после завершения middleware
//			defer cw.Close()
//			ow.Header().Set("Content-Encoding", "gzip")
//		}
//
//		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
//		contentEncoding := r.Header.Get("Content-Encoding")
//		//fmt.Println(contentEncoding) // для тнстов
//		sendsGzip := strings.Contains(contentEncoding, "gzip")
//		if sendsGzip {
//			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
//			cr, err := newCompressReader(r.Body)
//			if err != nil {
//				w.WriteHeader(http.StatusInternalServerError)
//				return
//			}
//			// меняем тело запроса на новое
//			r.Body = cr
//			defer cr.Close()
//		}
//
//		// передаём управление хендлеру
//		h(ow, r)
//	}
//}

func gzipBody(data []byte) ([]byte, error) {
	var b bytes.Buffer
	// создаём переменную w — в неё будут записываться входящие данные,
	// которые будут сжиматься и сохраняться в bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	// запись данных
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	// обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	// переменная b содержит сжатые данные
	return b.Bytes(), nil
}

func SendPostGzipJSON(client *http.Client, metrics *Metrics, url *string) error {

	body, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	compressed, err := gzipBody(body)
	if err != nil {
		return err
	}
	r := strings.NewReader(string(compressed))

	req, err := http.NewRequest(http.MethodPost, *url, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if !errors.Is(err, syscall.ECONNRESET) && !errors.Is(err, io.EOF) && !errors.Is(err, syscall.ECONNREFUSED) {
			return err
		}
	}
	req.Close = true

	return nil
}
