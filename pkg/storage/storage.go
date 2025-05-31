package storage

import (
	"Quotes1.0/internal/models"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type QuoteRepository interface {
	Add(input models.QuoteInput) (*models.Quote, error)
	GetAll() ([]models.Quote, error)
	GetRandom() (*models.Quote, error)
	GetByAuthor(author string) ([]models.Quote, error)
	DeleteByID(id int) (*models.Quote, error)
}

type MemoryStorage struct {
	quotes []models.Quote
	nextID int
	rnd    *rand.Rand
	mu     sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		quotes: make([]models.Quote, 0),
		nextID: 1,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *MemoryStorage) Add(input models.QuoteInput) (*models.Quote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	quote := models.Quote{
		ID:     m.nextID,
		Quote:  input.Quote,
		Author: input.Author,
	}
	m.nextID++
	m.quotes = append(m.quotes, quote)

	return &quote, nil
}

func (m *MemoryStorage) GetAll() ([]models.Quote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.quotes) == 0 {
		return nil, errors.New("no quotes found")
	}
	return m.quotes, nil
}

func (m *MemoryStorage) GetRandom() (*models.Quote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.quotes) == 0 {
		return nil, errors.New("no quotes available")
	}

	return &m.quotes[m.rnd.Intn(len(m.quotes))], nil
}

func (m *MemoryStorage) GetByAuthor(author string) ([]models.Quote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []models.Quote
	for _, q := range m.quotes {
		if q.Author == author {
			result = append(result, q)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no quotes found for author %s", author)
	}
	return result, nil
}

func (m *MemoryStorage) DeleteByID(id int) (*models.Quote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, q := range m.quotes {
		if q.ID == id {
			m.quotes = append(m.quotes[:i], m.quotes[i+1:]...)
			return &q, nil
		}
	}

	return nil, errors.New("quote not found")
}
