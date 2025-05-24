package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"quoter/internal/handlers"
	"quoter/internal/routes"
	"quoter/internal/storage"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	quoteStorage := storage.New()

	quoteHandler := &handlers.QuoteHandler{
		Storage: quoteStorage,
	}
	r := routes.NewRouter(quoteHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
