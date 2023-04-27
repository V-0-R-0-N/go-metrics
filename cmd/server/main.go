package main

import (
	"net/http"

	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handlers.BadRequest)
	mux.HandleFunc(`/update/`, handlers.UpdateMetrics)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
