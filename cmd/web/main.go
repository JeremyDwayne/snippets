package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", getSnippetView)
	mux.HandleFunc("GET /snippet/create", getSnippetCreate)
	mux.HandleFunc("POST /snippet/create", postSnippetCreate)

	port := os.Getenv("HTTP_LISTEN_ADDR")
	log.Printf("Starting server on %s", port)

	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}
