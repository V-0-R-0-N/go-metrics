package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/V-0-R-0-N/go-metrics.git/internal/checkers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/structs"
)

func BadRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func UpdateMetrics(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		splittedURI := strings.Split(strings.Trim(req.RequestURI, "/"), "/")

		if len(splittedURI) == 2 {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		if len(splittedURI) != 4 || !checkers.CheckMetricType(splittedURI[1]) || !checkers.CheckContentType(req) {

			BadRequest(res, req)
			return
		}

		metricType := splittedURI[1]
		metricName := splittedURI[2]
		metricData := splittedURI[3]
		if metricType == "gauge" {
			num, err := strconv.ParseFloat(metricData, 64)
			if err != nil {

				BadRequest(res, req)
				return
			}
			structs.Storage.GaugeData[metricName] = structs.Gauge(num)

		} else if metricType == "counter" {
			num, err := strconv.Atoi(metricData)
			if err != nil || num < 0 || structs.Storage.CounterData[metricName]+structs.Counter(num) < 0 {

				BadRequest(res, req)
				return
			}

			structs.Storage.CounterData[metricName] += structs.Counter(num)

		}

	} else {
		BadRequest(res, req)
	}

}
