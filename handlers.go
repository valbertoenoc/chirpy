package main

import (
	"net/http"
)

func handlerFileServer() http.Handler {
	filePathRoot := "."
	return http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
}
