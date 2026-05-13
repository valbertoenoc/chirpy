package main

import (
	"log"
	"net/http"

	"github.com/valbertoenoc/chirpy/internal/config"
	"github.com/valbertoenoc/chirpy/internal/handlers"
)

func main() {
	cfg := config.ApiConfig{}

	port := "8080"

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":" + "8080",
	}

	mux.Handle("/app/", cfg.MiddlewareMetricsInc(handlerFileServer()))
	mux.HandleFunc("GET /healthz", handlers.HandlerHealth)
	mux.HandleFunc("GET /metrics", cfg.HandlerMetrics)
	mux.HandleFunc("POST /reset", cfg.HandlerReset)

	log.Printf("Listening on port: %s", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server")
	}

}
