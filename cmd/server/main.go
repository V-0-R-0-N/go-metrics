package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"

	"github.com/V-0-R-0-N/go-metrics.git/internal/environ"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

func main() {
	addr := flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	flags.Server(&addr)
	flag.Parse()
	if err := environ.Server(&addr); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	//router.Use(middleware.Logger) // Для тестов

	st := storage.NewStorage()

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

	err := http.ListenAndServe(addr.String(), router)
	if err != nil {
		log.Fatal(err)
	}
}
