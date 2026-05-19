package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/valbertoenoc/chirpy/internal/database"
)

type apiConfig struct {
	db             *database.Queries
	fileserverHits atomic.Int32
	platform       string
	secretKey      string
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	secretKey := os.Getenv("SECRET_KEY")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(dbConn)

	cfg := apiConfig{
		db:        dbQueries,
		platform:  platform,
		secretKey: secretKey,
	}

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":" + "8080",
	}
	fsServer := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(fsServer))

	// /api namespace
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	// api namespace /chirps resource
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerListChirps)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.handlerGetChirp)
	// api namespace /users resource
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	// /admin namespace
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.HandlerReset)

	log.Printf("Listening on port: %s", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server")
	}

}
