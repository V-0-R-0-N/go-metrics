package checkers

import "net/http"

func CheckMetricType(name string) bool {
	types := make(map[string]bool)

	types["gauge"] = true
	types["counter"] = true

	if _, ok := types[name]; ok {
		return true
	}
	return false
}

func CheckContentType(req *http.Request) bool {
	return req.Header.Get("Content-Type") == "text/plain" || req.Header.Get("Content-Type") == ""
}
