package handlers

import (
	"Quotes1.0/internal/models"
	"Quotes1.0/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type QuoteHandler struct {
	repo storage.QuoteRepository
}

func NewQuoteHandler(repo storage.QuoteRepository) *QuoteHandler {
	return &QuoteHandler{repo: repo}
}

func (h *QuoteHandler) AddQuote(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var input models.QuoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.repo.Add(input)
	if err != nil {
		http.Error(w, "Failed to add quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quote)
}

func (h *QuoteHandler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {

	author := r.URL.Query().Get("author")

	if author != "" {
		quotes, err := h.repo.GetByAuthor(author)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(quotes)
		return
	}

	quotes, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get quotes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.repo.GetRandom()
	if err != nil {
		http.Error(w, "No quotes available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (h *QuoteHandler) GetQuotesByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author") // Получаем author из query-параметров
	if author == "" {
		http.Error(w, "Author parameter is required", http.StatusBadRequest)
		return
	}

	quotes, err := h.repo.GetByAuthor(author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *QuoteHandler) DeleteQuoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	deletedQuote, err := h.repo.DeleteByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(deletedQuote); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
