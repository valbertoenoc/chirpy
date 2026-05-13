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
	fsServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(fsServer))

	// /api namespace
	mux.HandleFunc("GET /api/healthz", handlers.HandlerHealth)

	// /admin namespace
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.HandlerReset)

	log.Printf("Listening on port: %s", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server")
	}

}
