package app

import (
	"database/sql"
	"net/http"

	"bookstore-api/internal/handlers"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	handler := handlers.NewHandler(db)

	mux.HandleFunc("/books", handler.BooksHandler)
	mux.HandleFunc("/books/", handler.BookHandler)

	mux.HandleFunc("/authors", handler.AuthorsHandler)
	mux.HandleFunc("/authors/", handler.AuthorHandler)

	mux.HandleFunc("/customers", handler.CustomersHandler)
	mux.HandleFunc("/customers/", handler.CustomerHandler)

	mux.HandleFunc("/purchases", handler.PurchasesHandler)

	return mux
}
