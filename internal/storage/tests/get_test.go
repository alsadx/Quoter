package tests

import (
	"strings"
	"testing"

	"quoter/internal/models"
	"quoter/internal/storage"
)

func TestGetAllQuotes(t *testing.T) {
	testStorage := storage.New()

	testStorage.AddQuote(models.Quote{Author: "Confucius", Text: "..."})
	testStorage.AddQuote(models.Quote{Author: "Lao Tzu", Text: "..."})

	quotes := testStorage.GetAllQuotes()
	if len(quotes) != 2 {
		t.Errorf("Expected 2 quotes, got %d", len(quotes))
	}

	if quotes[0].Author != "Confucius" || quotes[1].Author != "Lao Tzu" {
		t.Errorf("Expected quotes to have authors 'Confucius' and 'Lao Tzu'")
	}

	if quotes[0].Text != "..." || quotes[1].Text != "..." {
		t.Errorf("Expected quotes to have texts '...' and '...'")
	}

	if quotes[0].ID != 1 || quotes[1].ID != 2 {
		t.Errorf("Expected quotes to have IDs 1 and 2")
	}
}

func TestGetQuotesByAuthor(t *testing.T) {
	testStorage := storage.New()

	testStorage.AddQuote(models.Quote{Author: "Confucius", Text: "quote 1"})
	testStorage.AddQuote(models.Quote{Author: "Confucius", Text: "quote 2"})

	quotes := testStorage.GetQuotesByAuthor("Confucius")

	if len(quotes) != 2 {
		t.Errorf("Expected 2 quotes, got %d", len(quotes))
	}

	if quotes[0].Author != "Confucius" || quotes[1].Author != "Confucius" {
		t.Errorf("Expected quotes to have author 'Confucius'")
	}

	if quotes[0].Text != "quote 1" || quotes[1].Text != "quote 2" {
		t.Errorf("Expected quotes to have text 'quote 1' and 'quote 2'")
	}

	if quotes[0].ID != 1 || quotes[1].ID != 2 {
		t.Errorf("Expected quotes to have IDs 1 and 2")
	}
}

func TestGetRandomQuote(t *testing.T) {
	testStorage := storage.New()

	testStorage.AddQuote(models.Quote{Author: "Confucius", Text: "quote 1"})
	testStorage.AddQuote(models.Quote{Author: "Confucius", Text: "quote 2"})
	testStorage.AddQuote(models.Quote{Author: "Confucius", Text: "quote 3"})

	quote := testStorage.GetRandomQuote()

	if quote.Author != "Confucius" {
		t.Errorf("Expected quote to have author 'Confucius'")
	}

	if !strings.HasPrefix(quote.Text, "quote ") {
		t.Errorf("Expected quote to have text starting with 'quote '")
	}
}