package handlers

import (
	"Quotes1.0/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage реализует интерфейс QuoteRepository для тестов
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Add(input models.QuoteInput) (*models.Quote, error) {
	args := m.Called(input)
	return args.Get(0).(*models.Quote), args.Error(1)
}

func (m *MockStorage) GetAll() ([]models.Quote, error) {
	args := m.Called()
	return args.Get(0).([]models.Quote), args.Error(1)
}

func (m *MockStorage) GetRandom() (*models.Quote, error) {
	args := m.Called()
	return args.Get(0).(*models.Quote), args.Error(1)
}

func (m *MockStorage) GetByAuthor(author string) ([]models.Quote, error) {
	args := m.Called(author)
	return args.Get(0).([]models.Quote), args.Error(1)
}

func (m *MockStorage) DeleteByID(id int) (*models.Quote, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Quote), args.Error(1)
}

func TestAddQuote_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewQuoteHandler(mockStorage)

	quoteInput := models.QuoteInput{
		Author: "Test Author",
		Quote:  "Test Quote",
	}
	expectedQuote := &models.Quote{
		ID:     1,
		Author: "Test Author",
		Quote:  "Test Quote",
	}

	mockStorage.On("Add", quoteInput).Return(expectedQuote, nil)

	body, _ := json.Marshal(quoteInput)
	req := httptest.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddQuote(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.Quote
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, *expectedQuote, response)
	mockStorage.AssertExpectations(t)
}

func TestAddQuote_InvalidContentType(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewQuoteHandler(mockStorage)

	req := httptest.NewRequest("POST", "/quotes", nil)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	handler.AddQuote(w, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

func TestGetAllQuotes_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewQuoteHandler(mockStorage)

	expectedQuotes := []models.Quote{
		{ID: 1, Author: "Author1", Quote: "Quote1"},
		{ID: 2, Author: "Author2", Quote: "Quote2"},
	}

	mockStorage.On("GetAll").Return(expectedQuotes, nil)

	req := httptest.NewRequest("GET", "/quotes", nil)
	w := httptest.NewRecorder()

	handler.GetAllQuotes(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.Quote
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedQuotes, response)
	mockStorage.AssertExpectations(t)
}

func TestGetRandomQuote_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewQuoteHandler(mockStorage)

	expectedQuote := &models.Quote{
		ID:     1,
		Author: "Test Author",
		Quote:  "Test Quote",
	}

	mockStorage.On("GetRandom").Return(expectedQuote, nil)

	req := httptest.NewRequest("GET", "/quotes/random", nil)
	w := httptest.NewRecorder()

	handler.GetRandomQuote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.Quote
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, *expectedQuote, response)
	mockStorage.AssertExpectations(t)
}

func TestGetQuotesByAuthor_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewQuoteHandler(mockStorage)

	expectedQuotes := []models.Quote{
		{ID: 1, Author: "Author1", Quote: "Quote1"},
		{ID: 2, Author: "Author1", Quote: "Quote2"},
	}

	mockStorage.On("GetByAuthor", "Author1").Return(expectedQuotes, nil)

	req := httptest.NewRequest("GET", "/quotes?author=Author1", nil)
	w := httptest.NewRecorder()

	handler.GetAllQuotes(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.Quote
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedQuotes, response)
	mockStorage.AssertExpectations(t)
}

func TestDeleteQuoteByID_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewQuoteHandler(mockStorage)

	expectedQuote := &models.Quote{
		ID:     1,
		Author: "Test Author",
		Quote:  "Test Quote",
	}

	// Убедимся, что мок ожидает вызов с правильным ID
	mockStorage.On("DeleteByID", 1).Return(expectedQuote, nil)

	// Создаем запрос с правильным ID в URL
	req := httptest.NewRequest("DELETE", "/quotes/1", nil)

	// Добавляем переменные маршрута вручную, так как httptest.NewRequest
	// не интегрируется с gorilla/mux по умолчанию
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()

	handler.DeleteQuoteByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Quote
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, *expectedQuote, response)

	mockStorage.AssertExpectations(t)
}

func TestDeleteQuoteByID_InvalidID(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewQuoteHandler(mockStorage)

	req := httptest.NewRequest("DELETE", "/quotes/invalid", nil)
	w := httptest.NewRecorder()

	handler.DeleteQuoteByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
