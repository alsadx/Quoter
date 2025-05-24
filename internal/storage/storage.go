package storage

import (
	"math/rand"
	"quoter/internal/models"
	"slices"
	"sync"
)

type QuoteStorage struct {
	Mu       sync.RWMutex
	Quotes   []models.Quote
	ByID     map[int]models.Quote
	ByAuthor map[string][]models.Quote
	Counter  int
}

func New() *QuoteStorage {
	return &QuoteStorage{
		Quotes:   make([]models.Quote, 0),
		ByID:     make(map[int]models.Quote),
		ByAuthor: make(map[string][]models.Quote),
		Counter:  1,
	}
}

func (s *QuoteStorage) GetAllQuotes() []models.Quote {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.Quotes
}

func (s *QuoteStorage) GetRandomQuote() models.Quote {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	if len(s.Quotes) == 0 {
		return models.Quote{}
	}

	randIndex := rand.Intn(len(s.Quotes))
	return s.Quotes[randIndex]
}

func (s *QuoteStorage) GetQuotesByAuthor(author string) []models.Quote {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	if quotes, exists := s.ByAuthor[author]; exists {
		return quotes
	}

	return []models.Quote{}
}

func (s *QuoteStorage) AddQuote(quote models.Quote) models.Quote {

	s.Mu.Lock()
	defer s.Mu.Unlock()

	quote.ID = s.Counter
	s.Counter++

	s.Quotes = append(s.Quotes, quote)

	s.ByID[quote.ID] = quote

	s.ByAuthor[quote.Author] = append(s.ByAuthor[quote.Author], quote)

	return quote
}

func (s *QuoteStorage) DeleteQuote(id int) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	quote, exists := s.ByID[id]
	if !exists {
		return false
	}

	delete(s.ByID, id)

	authorQuotes := s.ByAuthor[quote.Author]
	for i, q := range authorQuotes {
		if q.ID == id {
			s.ByAuthor[quote.Author] = slices.Delete(authorQuotes, i, i+1)
			break
		}
	}

	for i, q := range s.Quotes {
		if q.ID == id {
			s.Quotes = slices.Delete(s.Quotes, i, i+1)
			break
		}
	}

	return true
}
