package main

import (
	"flag"
	"github.com/V-0-R-0-N/go-metrics.git/internal/environ"
	"github.com/V-0-R-0-N/go-metrics.git/internal/filer"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/middlware/compress"
	"github.com/V-0-R-0-N/go-metrics.git/internal/middlware/logger"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
)

//var Log, _ = zap.NewDevelopment()

func main() {
	//defer Log.Sync()

	Log, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	if Log != nil {
		defer Log.Sync()
	}
	sugar := logger.NewSugarLogger(Log)

	addr := flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	FileR := flags.NewFileR()
	flags.Server(&addr, FileR)
	flag.Parse()
	if err := environ.Server(&addr, FileR); err != nil {
		log.Fatal(err)
	}

	filer.FilerInit(FileR)
	if FileR.File.File != nil {
		defer filer.Close(&FileR.File)
	}
	sugar.Infow(
		"Server start",
		zap.String("address: ", addr.String()),
	)

	router := chi.NewRouter()
	//router.Use(middleware.Logger) // Для тестов
	//router.Use(middleware.Compress(5, "text/html", "application/json")) // для тестов
	st := storage.NewStorage()

	filer.StartRestore(st, FileR)
	handlerStorage := handlers.NewHandlerStorage(st)

	//router.Get("/", logger.WithLogging(handlerStorage.GetMetrics))
	//
	//router.HandleFunc("/update/*", logger.WithLogging(handlerStorage.UpdateMetrics))
	//
	//router.Get("/value/{type}/{name}", logger.WithLogging(handlerStorage.GetMetricsValue))
	//
	//router.Post("/update/", logger.WithLogging(handlerStorage.UpdateMetricJSON))
	//router.Post("/value/", logger.WithLogging(handlerStorage.GetMetricJSON))

	//router.Get("/",
	//	logger.WithLogging(compress.GzipMiddleware(handlerStorage.GetMetrics)))
	//
	//router.HandleFunc("/update/*",
	//	logger.WithLogging(compress.GzipMiddleware(handlerStorage.UpdateMetrics)))
	//
	//router.Get("/value/{type}/{name}",
	//	logger.WithLogging(compress.GzipMiddleware(handlerStorage.GetMetricsValue)))
	//
	//router.Post("/update/",
	//	logger.WithLogging(compress.GzipMiddleware(handlerStorage.UpdateMetricJSON)))
	//router.Post("/value/",
	//	logger.WithLogging(compress.GzipMiddleware(handlerStorage.GetMetricJSON)))

	// TODO обсудить эту реализацию с ментором

	router.Use(logger.WithLogging)
	router.Use(compress.GzipMiddleware)

	router.Get("/", handlerStorage.GetMetrics)

	router.HandleFunc("/update/*", handlerStorage.UpdateMetrics)

	router.Get("/value/{type}/{name}", handlerStorage.GetMetricsValue)

	router.Post("/update/", handlerStorage.UpdateMetricJSON)
	router.Post("/value/", handlerStorage.GetMetricJSON)

	err = http.ListenAndServe(addr.String(), router)
	if err != nil {
		log.Fatal(err)
	}
}
