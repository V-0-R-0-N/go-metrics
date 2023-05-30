package router

import (
	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/middlware/compress"
	"github.com/V-0-R-0-N/go-metrics.git/internal/middlware/logger"
	"github.com/go-chi/chi/v5"
)

func Router(logMiddlware logger.Middlware, handlerStorage *handlers.Handler) chi.Router {
	router := chi.NewRouter()
	//router.Use(middleware.Logger) // Для тестов
	router.Use(logMiddlware.WithLogging)
	router.Use(compress.GzipMiddleware)

	router.Get("/", handlerStorage.GetMetrics)

	router.HandleFunc("/update/*", handlerStorage.UpdateMetrics)

	router.Get("/value/{type}/{name}", handlerStorage.GetMetricsValue)

	router.Post("/update/", handlerStorage.UpdateMetricJSON)
	router.Post("/value/", handlerStorage.GetMetricJSON)

	return router
}
