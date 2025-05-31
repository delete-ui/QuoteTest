package main

import (
	"Quotes1.0/internal/handlers"
	"Quotes1.0/pkg/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	
	memoryRepository := storage.NewMemoryStorage()
	quoteHandler := handlers.NewQuoteHandler(memoryRepository)

	r := mux.NewRouter()
	r.HandleFunc("/quotes", quoteHandler.AddQuote).Methods("POST")
	r.HandleFunc("/quotes", quoteHandler.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", quoteHandler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes", quoteHandler.GetQuotesByAuthor).Methods("GET").Queries("author", "{author}")
	r.HandleFunc("/quotes/{id}", quoteHandler.DeleteQuoteByID).Methods("DELETE")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
