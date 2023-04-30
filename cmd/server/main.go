package main

import (
	"net/http"

	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handlers.BadRequest)

	st := storage.New()

	handlerUpdate := handlers.NewHandlerStorage(st)

	mux.HandleFunc(`/update/`, handlerUpdate.UpdateMetrics)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
