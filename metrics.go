package main

import (
	"fmt"
	"net/http"
)

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
		<html>
		<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
		</body>
		</html>
		`, c.fileserverHits.Load())))
}

func (c *apiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if c.platform != "DEV" {
		respondWithError(w, http.StatusForbidden, "forbidden operation.")
		return
	}

	err := c.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
