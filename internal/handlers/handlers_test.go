package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
}
