package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"

	ch "github.com/V-0-R-0-N/go-metrics.git/internal/checkers"
	st "github.com/V-0-R-0-N/go-metrics.git/internal/storage"
)

func BadRequest(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func updateValidator(splURI []string, req *http.Request) bool {
	if len(splURI) <= 2 ||
		len(splURI) != 4 ||
		!ch.CheckMetricType(splURI[1]) ||
		!ch.CheckContentType(req) {

		return true
	}
	return false
}

type handler struct {
	storage st.Storage
}

func NewHandlerStorage(storage st.Storage) *handler {
	return &handler{storage}
}

func (h *handler) UpdateMetrics(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		splURI := strings.Split(strings.Trim(req.RequestURI, "/"), "/")

		if len(splURI) == 2 && ch.CheckMetricType(splURI[1]) {
			res.WriteHeader(http.StatusNotFound)

			return
		}
		if updateValidator(splURI, req) {
			BadRequest(res, req)
			return
		}

		metricType := splURI[1]
		metricName := splURI[2]
		metricData := splURI[3]
		if metricType == "gauge" {
			num, err := strconv.ParseFloat(metricData, 64)
			if err != nil {
				BadRequest(res, req)
				return
			}

			h.storage.PutGauge(metricName, st.Float64ToGauge(num))
		} else if metricType == "counter" {
			num, err := strconv.Atoi(metricData)
			if err != nil || num < 0 || h.storage.GetCounter(metricName)+st.IntToCounter(num) < 0 {
				BadRequest(res, req)
				return
			}

			h.storage.PutCounter(metricName, st.IntToCounter(num))
		}
		//fmt.Println(h.storage.GetStorage()) // Для теста
		return
	}

	BadRequest(res, req)
}

func (h *handler) GetMetrics(res http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(res, fmt.Sprint(h.storage))
}

func (h *handler) GetMetricsValue(res http.ResponseWriter, req *http.Request) {

	t := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")
	if t != "gauge" && t != "counter" {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if t == "gauge" {
		if _, ok := h.storage.GetStorage().Gauge[name]; !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(res, fmt.Sprint(h.storage.GetStorage().Gauge[name]))
	}
	if t == "counter" {
		if _, ok := h.storage.GetStorage().Counter[name]; !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(res, fmt.Sprint(h.storage.GetStorage().Counter[name]))
	}
}
