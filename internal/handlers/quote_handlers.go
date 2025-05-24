package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"quoter/internal/models"
	"strconv"

	"github.com/gorilla/mux"
)

type Storage interface {
	GetAllQuotes() []models.Quote
	GetRandomQuote() models.Quote
	GetQuotesByAuthor(author string) []models.Quote
	AddQuote(quote models.Quote) models.Quote
	DeleteQuote(id int) bool
}

type QuoteHandler struct {
	Storage Storage
}

func (h *QuoteHandler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	quotes := h.Storage.GetAllQuotes()
	json.NewEncoder(w).Encode(quotes)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote := h.Storage.GetRandomQuote()
	if quote.ID == 0 {
		http.Error(w, "No quotes available", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(quote)
}

func (h *QuoteHandler) GetQuotesByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	quotes := h.Storage.GetQuotesByAuthor(author)

	json.NewEncoder(w).Encode(quotes)
}

func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var quote models.Quote
	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := quote.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newQuote := h.Storage.AddQuote(quote)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newQuote)
}

func (h *QuoteHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	if !h.Storage.DeleteQuote(id) {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Quote with ID %d deleted", id)
}
