package config

import (
	"net/http"
	"strconv"
	"sync/atomic"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
}

func (c *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	responseBody := "Hits: " + strconv.FormatInt(int64(c.fileserverHits.Load()), 10)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBody))
}

func (c *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
