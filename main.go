package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nof20/shortener/shortener"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", shortener.RootGetHandler).Methods("GET")
	r.HandleFunc("/", shortener.RootPostHandler).Methods("POST")
	r.HandleFunc("/{key}", shortener.URLHandler)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
