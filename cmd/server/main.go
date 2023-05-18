package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/V-0-R-0-N/go-metrics.git/internal/environ"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/logger"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

//var Log, _ = zap.NewDevelopment()

func main() {
	//defer Log.Sync()

	Log, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer Log.Sync()
	sugar := logger.NewSugarLogger(Log)

	addr := flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	sugar.Infow(
		"Server start",
		zap.String("address: ", addr.String()),
	)
	flags.Server(&addr)
	flag.Parse()
	if err := environ.Server(&addr); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	//router.Use(middleware.Logger) // Для тестов

	st := storage.NewStorage()

	handlerStorage := handlers.NewHandlerStorage(st)

	router.Get("/", logger.WithLogging(handlerStorage.GetMetrics))

	router.HandleFunc("/update/*", logger.WithLogging(handlerStorage.UpdateMetrics))

	router.Get("/value/{type}/{name}", logger.WithLogging(handlerStorage.GetMetricsValue))

	err = http.ListenAndServe(addr.String(), router)
	if err != nil {
		log.Fatal(err)
	}
}
