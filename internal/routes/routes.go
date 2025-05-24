package routes

import (
    "github.com/gorilla/mux"
    "quoter/internal/handlers"
)

func NewRouter(handler *handlers.QuoteHandler) *mux.Router {
    r := mux.NewRouter()

	r.HandleFunc("/quotes", handler.GetQuotesByAuthor).Queries("author", "{author}").Methods("GET")
    r.HandleFunc("/quotes", handler.GetAllQuotes).Methods("GET")
    r.HandleFunc("/quotes/random", handler.GetRandomQuote).Methods("GET")
    r.HandleFunc("/quotes", handler.CreateQuote).Methods("POST")
    r.HandleFunc("/quotes/{id}", handler.DeleteQuote).Methods("DELETE")

    return r
}