package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

//type Logger interface {
//	Infoln(args ...interface{})
//}

var sugar *zap.SugaredLogger

//type LoggerMiddlware struct {
//	sugar Logger
//}

func NewSugarLogger(Log *zap.Logger) *zap.SugaredLogger {
	//return
	//return &LoggerMiddlware{
	//	sugar: log,
	//}
	sugar = Log.Sugar()
	return sugar
}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// функция Now() возвращает текущее время
		start := time.Now()

		// эндпоинт
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		// данные ответа
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		// точка, где выполняется хендлер
		h.ServeHTTP(&lw, r) // обслуживание оригинального запроса

		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
		duration := time.Since(start)

		// отправляем сведения о запросе в zap
		sugar.Infoln(
			//"HTTP request",
			"URI", uri,
			"method", method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	})
}

//func WithLogging(h http.HandlerFunc) http.HandlerFunc {
//	logFn := func(w http.ResponseWriter, r *http.Request) {
//		// функция Now() возвращает текущее время
//		start := time.Now()
//
//		// эндпоинт
//		uri := r.RequestURI
//		// метод запроса
//		method := r.Method
//
//		// данные ответа
//		responseData := &responseData{
//			status: 0,
//			size:   0,
//		}
//		lw := loggingResponseWriter{
//			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
//			responseData:   responseData,
//		}
//		// точка, где выполняется хендлер
//		h(&lw, r) // обслуживание оригинального запроса
//
//		// Since возвращает разницу во времени между start
//		// и моментом вызова Since. Таким образом можно посчитать
//		// время выполнения запроса.
//		duration := time.Since(start)
//
//		// отправляем сведения о запросе в zap
//		sugar.Infoln(
//			//"HTTP request",
//			"URI", uri,
//			"method", method,
//			"status", responseData.status,
//			"duration", duration,
//			"size", responseData.size,
//		)
//
//	}
//	// возвращаем функционально расширенный хендлер
//	return http.HandlerFunc(logFn)
//}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
