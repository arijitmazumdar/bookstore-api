package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"bookstore-api/internal/db"
	"bookstore-api/internal/models"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	tmpFile, err := os.CreateTemp("", "bookstore_test_*.db")
	if err != nil {
		t.Fatalf("create temp db: %v", err)
	}
	dbPath := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("close temp file: %v", err)
	}

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := db.RunMigrations(database); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	return database, func() {
		if err := database.Close(); err != nil {
			t.Fatalf("close database: %v", err)
		}
		if err := os.Remove(dbPath); err != nil {
			t.Fatalf("remove db file: %v", err)
		}
	}
}

func TestAuthorLifecycle(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	h := NewHandler(database)

	newAuthor := models.Author{Name: "Jane Doe"}
	body, _ := json.Marshal(newAuthor)
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.AuthorsHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var created models.Author
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if created.ID == 0 || created.Name != newAuthor.Name {
		t.Fatalf("unexpected author created: %+v", created)
	}
	if created.Category != authorCategoryRegular {
		t.Fatalf("expected created author category %q, got %+v", authorCategoryRegular, created)
	}

	updateBody, _ := json.Marshal(models.Author{Name: "Jane Updated"})
	updateReq := httptest.NewRequest(http.MethodPut, "/authors/"+itoa(created.ID), bytes.NewReader(updateBody))
	updateW := httptest.NewRecorder()
	h.AuthorHandler(updateW, updateReq)

	if updateW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", updateW.Code, updateW.Body.String())
	}

	var updated models.Author
	if err := json.NewDecoder(updateW.Body).Decode(&updated); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if updated.Name != "Jane Updated" || updated.Category != authorCategoryRegular {
		t.Fatalf("unexpected author updated: %+v", updated)
	}
}

func TestAuthorsListAndDetailReturnRegularCategoryWithoutSales(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	h := NewHandler(database)

	authorID := insertAuthor(t, database, "Author Without Sales")

	req := httptest.NewRequest(http.MethodGet, "/authors", nil)
	w := httptest.NewRecorder()
	h.AuthorsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for list, got %d: %s", w.Code, w.Body.String())
	}

	var authors []models.Author
	if err := json.NewDecoder(w.Body).Decode(&authors); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(authors) != 1 {
		t.Fatalf("expected 1 author in list, got %d", len(authors))
	}
	if authors[0].ID != authorID || authors[0].Category != authorCategoryRegular {
		t.Fatalf("unexpected authors list entry: %+v", authors[0])
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/authors/"+itoa(authorID), nil)
	detailW := httptest.NewRecorder()
	h.AuthorHandler(detailW, detailReq)

	if detailW.Code != http.StatusOK {
		t.Fatalf("expected 200 for detail, got %d: %s", detailW.Code, detailW.Body.String())
	}

	var author models.Author
	if err := json.NewDecoder(detailW.Body).Decode(&author); err != nil {
		t.Fatalf("decode detail response: %v", err)
	}
	if author.ID != authorID || author.Category != authorCategoryRegular {
		t.Fatalf("unexpected author detail: %+v", author)
	}
}

func TestAuthorCategoryThresholdAndAggregation(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	h := NewHandler(database)

	authorID := insertAuthor(t, database, "Threshold Author")
	firstBookID := insertBook(t, database, authorID, "First Book")
	secondBookID := insertBook(t, database, authorID, "Second Book")
	customerID := insertCustomer(t, database, "Buyer")

	insertPurchases(t, database, customerID, firstBookID, 250)
	insertPurchases(t, database, customerID, secondBookID, 250)

	author := fetchAuthorDetail(t, h, authorID)
	if author.Category != authorCategoryRegular {
		t.Fatalf("expected regular category at 500 sales, got %+v", author)
	}

	insertPurchases(t, database, customerID, secondBookID, 1)

	author = fetchAuthorDetail(t, h, authorID)
	if author.Category != authorCategoryPremium {
		t.Fatalf("expected premium category at 501 sales, got %+v", author)
	}
}

func TestAuthorCategoryRecalculatesFromCurrentPurchaseData(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	h := NewHandler(database)

	authorID := insertAuthor(t, database, "Recalc Author")
	bookID := insertBook(t, database, authorID, "Mutable Book")
	customerID := insertCustomer(t, database, "Buyer")

	insertPurchases(t, database, customerID, bookID, 501)

	author := fetchAuthorDetail(t, h, authorID)
	if author.Category != authorCategoryPremium {
		t.Fatalf("expected premium category before deletion, got %+v", author)
	}

	if _, err := database.Exec(
		"DELETE FROM customer_book_purchase WHERE id IN (SELECT id FROM customer_book_purchase WHERE book_id = ? ORDER BY id LIMIT 2)",
		bookID,
	); err != nil {
		t.Fatalf("delete purchases: %v", err)
	}

	author = fetchAuthorDetail(t, h, authorID)
	if author.Category != authorCategoryRegular {
		t.Fatalf("expected regular category after deletion, got %+v", author)
	}
}

func TestFailedPurchaseCreationDoesNotChangeAuthorCategory(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	h := NewHandler(database)

	authorID := insertAuthor(t, database, "Stable Author")
	customerID := insertCustomer(t, database, "Buyer")

	reqBody, _ := json.Marshal(models.Purchase{
		CustomerID:    customerID,
		BookID:        999999,
		PurchaseDate:  "2026-06-21T00:00:00Z",
		PurchasePrice: 19.99,
	})
	req := httptest.NewRequest(http.MethodPost, "/purchases", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	h.PurchasesHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for invalid purchase creation, got %d: %s", w.Code, w.Body.String())
	}

	author := fetchAuthorDetail(t, h, authorID)
	if author.Category != authorCategoryRegular {
		t.Fatalf("expected failed purchase to keep regular category, got %+v", author)
	}
}

func fetchAuthorDetail(t *testing.T, h *Handler, authorID int64) models.Author {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/authors/"+itoa(authorID), nil)
	w := httptest.NewRecorder()
	h.AuthorHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 fetching author %d, got %d: %s", authorID, w.Code, w.Body.String())
	}

	var author models.Author
	if err := json.NewDecoder(w.Body).Decode(&author); err != nil {
		t.Fatalf("decode author detail: %v", err)
	}

	return author
}

func insertAuthor(t *testing.T, database *sql.DB, name string) int64 {
	t.Helper()

	result, err := database.Exec("INSERT INTO authors (name) VALUES (?)", name)
	if err != nil {
		t.Fatalf("insert author: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("author last insert id: %v", err)
	}

	return id
}

func insertBook(t *testing.T, database *sql.DB, authorID int64, name string) int64 {
	t.Helper()

	result, err := database.Exec(
		"INSERT INTO books (name, print_type, paperback, author_id, publisher_house, copy_available) VALUES (?, ?, ?, ?, ?, ?)",
		name,
		"paperback",
		true,
		authorID,
		"Acme Publishing",
		true,
	)
	if err != nil {
		t.Fatalf("insert book: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("book last insert id: %v", err)
	}

	return id
}

func insertCustomer(t *testing.T, database *sql.DB, name string) int64 {
	t.Helper()

	result, err := database.Exec("INSERT INTO customers (name) VALUES (?)", name)
	if err != nil {
		t.Fatalf("insert customer: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("customer last insert id: %v", err)
	}

	return id
}

func insertPurchases(t *testing.T, database *sql.DB, customerID, bookID int64, count int) {
	t.Helper()

	for i := 0; i < count; i++ {
		if _, err := database.Exec(
			"INSERT INTO customer_book_purchase (customer_id, book_id, purchase_date, purchase_price) VALUES (?, ?, ?, ?)",
			customerID,
			bookID,
			"2026-06-21T00:00:00Z",
			19.99,
		); err != nil {
			t.Fatalf("insert purchase %d/%d: %v", i+1, count, err)
		}
	}
}

func itoa(id int64) string {
	return strconv.FormatInt(id, 10)
}
