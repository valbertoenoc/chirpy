package main

import "net/http"

// @Summary Health check
// @Description Returns the health status of the API
// @Tags health
// @Accept plain
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /api/healthz [get]
func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
