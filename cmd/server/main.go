package main

import (
	"flag"
	"fmt"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

var addr = flags.NetAddress{
	Host: "localhost",
	Port: 8080,
}

func init() {

	flags.Server(&addr)
}

func main() {

	flag.Parse()
	router := chi.NewRouter()
	//router.Use(middleware.Logger) // Для тестов

	st := storage.New()

	handlerStorage := handlers.NewHandlerStorage(st)

	router.Get("/", handlerStorage.GetMetrics)

	// TODO обсудить с ментором

	//router.Route("/update", func(r chi.Router) {
	//	router.Post("/", handlers.BadRequest)
	//	router.Route("/{type}", func(r chi.Router) {
	//		router.Post("/", handlers.BadRequest)
	//		router.Route("/{name}", func(r chi.Router) {
	//			router.Post("/", handlers.BadRequest)
	//			router.Post("/{data}", handlerStorage.UpdateMetrics)
	//		})
	//	})
	//})
	router.HandleFunc("/update/*", handlerStorage.UpdateMetrics)

	router.Get("/value/{type}/{name}", handlerStorage.GetMetricsValue)
	err := http.ListenAndServe(fmt.Sprintf("%s", addr), router)
	if err != nil {
		panic(err)
	}
}
