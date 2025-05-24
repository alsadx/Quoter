package tests

import (
	"testing"

	"quoter/internal/models"
	"quoter/internal/storage"
)

func TestInMemoryStore_AddQuote(t *testing.T) {
	testStorage := storage.New()

	quote := models.Quote{
		Author: "Confucius",
		Text:   "...",
	}

	newQuote := testStorage.AddQuote(quote)

	if newQuote.ID == 0 {
		t.Errorf("Expected non-zero ID")
	}

	if newQuote.Author != quote.Author || newQuote.Text != quote.Text {
		t.Errorf("Quote data mismatch")
	}
}

func TestInMemoryStore_DeleteQuote(t *testing.T) {
	testStorage := storage.New()

	q1 := testStorage.AddQuote(models.Quote{Author: "Confucius", Text: "..."})
	q2 := testStorage.AddQuote(models.Quote{Author: "Lao Tzu", Text: "..."})

	if !testStorage.DeleteQuote(q1.ID) {
		t.Errorf("Expected true from DeleteQuote")
	}

	if testStorage.DeleteQuote(999) {
		t.Errorf("Expected false when deleting non-existent ID")
	}

	quotes := testStorage.GetAllQuotes()
	if len(quotes) != 1 {
		t.Errorf("Expected 1 quote after deletion, got %d", len(quotes))
	}

	if quotes[0].ID != q2.ID {
		t.Errorf("Expected remaining quote to be second one")
	}
}
