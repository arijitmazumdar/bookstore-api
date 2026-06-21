package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bookstore-api/internal/models"
)

const (
	authorCategoryRegular = "regular"
	authorCategoryPremium = "premium"
)

func (h *Handler) AuthorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAuthors(w, r)
	case http.MethodPost:
		h.createAuthor(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) AuthorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAuthor(w, r, id)
	case http.MethodPut:
		h.updateAuthor(w, r, id)
	case http.MethodDelete:
		h.deleteAuthor(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listAuthors(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT
			a.id,
			a.name,
			COUNT(p.id) AS sold_copies
		FROM authors a
		LEFT JOIN books b ON b.author_id = a.id
		LEFT JOIN customer_book_purchase p ON p.book_id = b.id
		GROUP BY a.id, a.name
		ORDER BY a.id
	`)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	authors := []models.Author{}
	for rows.Next() {
		var author models.Author
		if err := rows.Scan(&author.ID, &author.Name, &author.SoldCopies); err != nil {
			h.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		author.Category = deriveAuthorCategory(author.SoldCopies)
		authors = append(authors, author)
	}

	if err := rows.Err(); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, authors)
}

func (h *Handler) getAuthor(w http.ResponseWriter, r *http.Request, id int64) {
	author, err := h.authorByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "author not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, author)
}

func (h *Handler) createAuthor(w http.ResponseWriter, r *http.Request) {
	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	result, err := h.db.Exec("INSERT INTO authors (name) VALUES (?)", author.Name)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	author.ID, _ = result.LastInsertId()
	createdAuthor, err := h.authorByID(author.ID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, createdAuthor)
}

func (h *Handler) updateAuthor(w http.ResponseWriter, r *http.Request, id int64) {
	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	author.ID = id
	result, err := h.db.Exec("UPDATE authors SET name = ? WHERE id = ?", author.Name, id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if rowsAffected == 0 {
		h.writeError(w, http.StatusNotFound, "author not found")
		return
	}

	updatedAuthor, err := h.authorByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "author not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, updatedAuthor)
}

func (h *Handler) deleteAuthor(w http.ResponseWriter, r *http.Request, id int64) {
	_, err := h.db.Exec("DELETE FROM authors WHERE id = ?", id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) authorByID(id int64) (models.Author, error) {
	var author models.Author
	row := h.db.QueryRow(`
		SELECT
			a.id,
			a.name,
			COUNT(p.id) AS sold_copies
		FROM authors a
		LEFT JOIN books b ON b.author_id = a.id
		LEFT JOIN customer_book_purchase p ON p.book_id = b.id
		WHERE a.id = ?
		GROUP BY a.id, a.name
	`, id)

	if err := row.Scan(&author.ID, &author.Name, &author.SoldCopies); err != nil {
		return models.Author{}, err
	}

	author.Category = deriveAuthorCategory(author.SoldCopies)
	return author, nil
}

func deriveAuthorCategory(soldCopies int64) string {
	if soldCopies > 500 {
		return authorCategoryPremium
	}

	return authorCategoryRegular
}
