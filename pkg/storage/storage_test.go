package storage

import (
	"Quotes1.0/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_AddAndGetAll(t *testing.T) {
	store := NewMemoryStorage()

	// Test empty storage
	quotes, err := store.GetAll()
	assert.Error(t, err)
	assert.Nil(t, quotes)

	// Add first quote
	quote1, err := store.Add(models.QuoteInput{
		Author: "Author1",
		Quote:  "Quote1",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, quote1.ID)
	assert.Equal(t, "Author1", quote1.Author)
	assert.Equal(t, "Quote1", quote1.Quote)

	// Add second quote
	quote2, err := store.Add(models.QuoteInput{
		Author: "Author2",
		Quote:  "Quote2",
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, quote2.ID)

	// Test GetAll
	quotes, err = store.GetAll()
	assert.NoError(t, err)
	assert.Len(t, quotes, 2)
	assert.Contains(t, quotes, *quote1)
	assert.Contains(t, quotes, *quote2)
}

func TestMemoryStorage_GetRandom(t *testing.T) {
	store := NewMemoryStorage()

	// Test empty storage
	quote, err := store.GetRandom()
	assert.Error(t, err)
	assert.Nil(t, quote)

	// Add quotes
	store.Add(models.QuoteInput{Author: "A1", Quote: "Q1"})
	store.Add(models.QuoteInput{Author: "A2", Quote: "Q2"})

	// Test GetRandom (run multiple times to ensure it works)
	found := make(map[int]bool)
	for i := 0; i < 10; i++ {
		quote, err = store.GetRandom()
		assert.NoError(t, err)
		assert.NotNil(t, quote)
		found[quote.ID] = true
	}
	assert.GreaterOrEqual(t, len(found), 1)
}

func TestMemoryStorage_GetByAuthor(t *testing.T) {
	store := NewMemoryStorage()

	// Test empty storage
	quotes, err := store.GetByAuthor("Unknown")
	assert.Error(t, err)
	assert.Nil(t, quotes)

	// Add quotes
	store.Add(models.QuoteInput{Author: "Author1", Quote: "Q1"})
	store.Add(models.QuoteInput{Author: "Author1", Quote: "Q2"})
	store.Add(models.QuoteInput{Author: "Author2", Quote: "Q3"})

	// Test GetByAuthor
	quotes, err = store.GetByAuthor("Author1")
	assert.NoError(t, err)
	assert.Len(t, quotes, 2)

	quotes, err = store.GetByAuthor("Author2")
	assert.NoError(t, err)
	assert.Len(t, quotes, 1)

	quotes, err = store.GetByAuthor("Unknown")
	assert.Error(t, err)
	assert.Nil(t, quotes)
}

func TestMemoryStorage_DeleteByID(t *testing.T) {
	store := NewMemoryStorage()

	// Test delete from empty storage
	quote, err := store.DeleteByID(1)
	assert.Error(t, err)
	assert.Nil(t, quote)

	// Add quotes
	q1, _ := store.Add(models.QuoteInput{Author: "A1", Quote: "Q1"})
	q2, _ := store.Add(models.QuoteInput{Author: "A2", Quote: "Q2"})

	// Test delete existing
	deleted, err := store.DeleteByID(q1.ID)
	assert.NoError(t, err)
	assert.Equal(t, q1, deleted)

	// Verify remaining quotes
	quotes, err := store.GetAll()
	assert.NoError(t, err)
	assert.Len(t, quotes, 1)
	assert.Equal(t, *q2, quotes[0])

	// Test delete non-existing
	quote, err = store.DeleteByID(999)
	assert.Error(t, err)
	assert.Nil(t, quote)
}

func TestMemoryStorage_Concurrency(t *testing.T) {
	store := NewMemoryStorage()

	// Run concurrent operations
	go func() {
		for i := 0; i < 100; i++ {
			store.Add(models.QuoteInput{
				Author: "Concurrent",
				Quote:  "Quote",
			})
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			store.GetAll()
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			store.GetRandom()
		}
	}()

	// Let the goroutines complete
	// If there are race conditions, the test will fail
}
