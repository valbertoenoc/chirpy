package main

import (
	"log"
	"net/http"
)

func main() {

	serve := http.NewServeMux()
	server := http.Server{
		Handler: serve,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server")
	}

}