package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

func main() {

	//mux := http.NewServeMux()

	router := chi.NewRouter()
	router.Use(middleware.Logger) // Для тестов

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
	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
