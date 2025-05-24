package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"quoter/internal/models"
	"quoter/internal/storage"

	"github.com/gorilla/mux"
)

func setupRouter() http.Handler {
	quoteStorage := storage.New()
	quoteHandler := QuoteHandler{Storage: quoteStorage}
	r := mux.NewRouter()
	r.HandleFunc("/quotes", quoteHandler.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", quoteHandler.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", quoteHandler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", quoteHandler.DeleteQuote).Methods("DELETE")

	return r
}

func TestAddDeleteQuote_HappyPath(t *testing.T) {
	router := setupRouter()

	author := "Confucius"
	quote := "Life is simple, but we insist on making it complicated."

	body := bytes.NewBufferString(`{"author":"` + author + `","quote":"` + quote + `"}`)

	req, _ := http.NewRequest("POST", "/quotes", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", rec.Code)
	}

	var response models.Quote
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Author != author || response.Text != quote {
		t.Errorf("Response data mismatch, got author %v and text %v", response.Author, response.Text)
	}

	if response.ID != 1 {
		t.Errorf("Expected ID 1; got %d", response.ID)
	}

	delReq := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(response.ID), nil)
	delRec := httptest.NewRecorder()
	router.ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusNoContent {
		t.Errorf("Expected status OK; got %v", delRec.Code)
	}
}

func TestAddQuote_EmptyParams(t *testing.T) {
	router := setupRouter()

	body := bytes.NewBufferString(`{"author":"   ","quote":"Life is simple, but we insist on making it complicated."}`)

	req, _ := http.NewRequest("POST", "/quotes", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request; got %v", rec.Code)
	}

	body = bytes.NewBufferString(`{"author":"Confucius","quote":""}`)

	req2, _ := http.NewRequest("POST", "/quotes", body)
	req2.Header.Set("Content-Type", "application/json")

	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request; got %v", rec2.Code)
	}
}

func TestAddQuote_InvalidParams(t *testing.T) {
	router := setupRouter()

	body := bytes.NewBufferString(`{"author":"Confucius","invalid_param":"Life is simple, but we insist on making it complicated."}`)

	req, _ := http.NewRequest("POST", "/quotes", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request; got %v", rec.Code)
	}
}

func TestDeleteQuote_InvalidID(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/quotes/invalid_id", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status Not Found; got %v", rec.Code)
	}
}

func TestDeleteQuote_NonexistentID(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/quotes/"+strconv.Itoa(999), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status Not Found; got %v", rec.Code)
	}
}

func TestGetAllQuotes_HappyPath(t *testing.T) {
	router := setupRouter()

	author1 := "Confucius"
	quote1 := "Life is simple, but we insist on making it complicated."

	author2 := "Jimmy Carr"
	quote2 := "Everyone is jealous of what you've got, no one is jealous of how you got it."

	body := bytes.NewBufferString(`{"author":"` + author1 + `","quote":"` + quote1 + `"}`)

	req, _ := http.NewRequest("POST", "/quotes", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", rec.Code)
	}

	body = bytes.NewBufferString(`{"author":"` + author2 + `","quote":"` + quote2 + `"}`)

	req2, _ := http.NewRequest("POST", "/quotes", body)
	req2.Header.Set("Content-Type", "application/json")

	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", rec2.Code)
	}

	getReq, _ := http.NewRequest("GET", "/quotes", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", getRec.Code)
	}

	var quotes []models.Quote
	if err := json.Unmarshal(getRec.Body.Bytes(), &quotes); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(quotes) != 2 {
		t.Errorf("Expected 2 quotes; got %d", len(quotes))
	}

	if quotes[0].Author != author1 || quotes[0].Text != quote1 {
		t.Errorf("Response data mismatch, got author1 %v and text1 %v", quotes[0].Author, quotes[0].Text)
	}

	if quotes[1].Author != author2 || quotes[1].Text != quote2 {
		t.Errorf("Response data mismatch, got author2 %v and text2 %v", quotes[1].Author, quotes[1].Text)
	}

	if quotes[0].ID != 1 || quotes[1].ID != 2 {
		t.Errorf("Expected quotes to have IDs 1 and 2, got %d and %d", quotes[0].ID, quotes[1].ID)
	}
}

func TestGetAllQuotes_Empty(t *testing.T) {
	router := setupRouter()

	getReq, _ := http.NewRequest("GET", "/quotes", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", getRec.Code)
	}

	var quotes []models.Quote
	if err := json.Unmarshal(getRec.Body.Bytes(), &quotes); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(quotes) != 0 {
		t.Errorf("Expected 0 quotes; got %d", len(quotes))
	}
}

func TestGetQuotesByAuthor_HappyPath(t *testing.T) {
	router := setupRouter()

	author := "Confucius"
	quote := "Life is simple, but we insist on making it complicated."

	body := bytes.NewBufferString(`{"author":"` + author + `","quote":"` + quote + `"}`)

	req, _ := http.NewRequest("POST", "/quotes", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", rec.Code)
	}

	getReq, _ := http.NewRequest("GET", "/quotes?author="+author, nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", getRec.Code)
	}

	var quotes []models.Quote
	if err := json.Unmarshal(getRec.Body.Bytes(), &quotes); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(quotes) != 1 {
		t.Errorf("Expected 1 quote; got %d", len(quotes))
	}

	if quotes[0].Author != author || quotes[0].Text != quote {
		t.Errorf("Response data mismatch, got author %v and text %v", quotes[0].Author, quotes[0].Text)
	}
}

func TestGetQuotesByAuthor_NotFound(t *testing.T) {
	router := setupRouter()

	getReq, _ := http.NewRequest("GET", "/quotes?author=NonexistentAuthor", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", getRec.Code)
	}

	var quotes []models.Quote
	if err := json.Unmarshal(getRec.Body.Bytes(), &quotes); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(quotes) != 0 {
		t.Errorf("Expected 0 quotes; got %d", len(quotes))
	}

	fmt.Println(quotes)
}

func TestGetQuotesByAuthor_EmptyAuthor(t *testing.T) {
	router := setupRouter()

	getReq, _ := http.NewRequest("GET", "/quotes?author=", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", getRec.Code)
	}

	var quotes []models.Quote
	if err := json.Unmarshal(getRec.Body.Bytes(), &quotes); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(quotes) != 0 {
		t.Errorf("Expected 0 quotes; got %d", len(quotes))
	}
}

func TestGetRandomQuote_HappyPath(t *testing.T) {
	router := setupRouter()

	author1 := "Confucius"
	quote1 := "Life is simple, but we insist on making it complicated."

	author2 := "Jimmy Carr"
	quote2 := "Everyone is jealous of what you've got, no one is jealous of how you got it"

	body := bytes.NewBufferString(`{"author":"` + author1 + `","quote":"` + quote1 + `"}`)

	req, _ := http.NewRequest("POST", "/quotes", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", rec.Code)
	}

	body = bytes.NewBufferString(`{"author":"` + author2 + `","quote":"` + quote2 + `"}`)

	req2, _ := http.NewRequest("POST", "/quotes", body)
	req2.Header.Set("Content-Type", "application/json")

	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", rec2.Code)
	}

	getReq, _ := http.NewRequest("GET", "/quotes/random", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", getRec.Code)
	}

	var quote models.Quote
	if err := json.Unmarshal(getRec.Body.Bytes(), &quote); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if quote.Author != author1 && quote.Author != author2 {
		t.Errorf("Got nonexistent author %v", quote.Author)
	}

	if quote.Text != quote1 && quote.Text != quote2 {
		t.Errorf("Got nonexistent quote %v", quote.Text)
	}

	if quote.ID != 1 && quote.ID != 2 {
		t.Errorf("Expected quotes to have IDs 1 and 2")
	}
}