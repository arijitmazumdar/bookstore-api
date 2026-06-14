package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"bookstore-api/internal/models"
)

func (h *Handler) BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listBooks(w, r)
	case http.MethodPost:
		h.createBook(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) BookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getBook(w, r, id)
	case http.MethodPut:
		h.updateBook(w, r, id)
	case http.MethodDelete:
		h.deleteBook(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, print_type, paperback, author_id, publisher_house, copy_available FROM books")
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	books := []models.Book{}
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Name, &book.PrintType, &book.Paperback, &book.AuthorID, &book.PublisherHouse, &book.CopyAvailable); err != nil {
			h.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		books = append(books, book)
	}

	h.writeJSON(w, http.StatusOK, books)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request, id int64) {
	var book models.Book
	row := h.db.QueryRow("SELECT id, name, print_type, paperback, author_id, publisher_house, copy_available FROM books WHERE id = ?", id)
	if err := row.Scan(&book.ID, &book.Name, &book.PrintType, &book.Paperback, &book.AuthorID, &book.PublisherHouse, &book.CopyAvailable); err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "book not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, book)
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	result, err := h.db.Exec("INSERT INTO books (name, print_type, paperback, author_id, publisher_house, copy_available) VALUES (?, ?, ?, ?, ?, ?)", book.Name, book.PrintType, book.Paperback, book.AuthorID, book.PublisherHouse, book.CopyAvailable)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	book.ID, _ = result.LastInsertId()
	h.writeJSON(w, http.StatusCreated, book)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request, id int64) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	book.ID = id
	_, err := h.db.Exec("UPDATE books SET name = ?, print_type = ?, paperback = ?, author_id = ?, publisher_house = ?, copy_available = ? WHERE id = ?", book.Name, book.PrintType, book.Paperback, book.AuthorID, book.PublisherHouse, book.CopyAvailable, id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, book)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request, id int64) {
	_, err := h.db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDFromPath(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, http.ErrMissingFile
	}
	return strconv.ParseInt(parts[len(parts)-1], 10, 64)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
