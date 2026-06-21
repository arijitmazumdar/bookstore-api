package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bookstore-api/internal/db"
	"bookstore-api/internal/models"
)

func setupIntegrationServer(t *testing.T) (*httptest.Server, *sql.DB, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "bookstore_integration_*.db")
	if err != nil {
		t.Fatal(err)
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

	router := NewRouter(database)
	server := httptest.NewServer(router)

	return server, database, func() {
		server.Close()
		if err := database.Close(); err != nil {
			t.Fatalf("close database: %v", err)
		}
		if err := os.Remove(dbPath); err != nil {
			t.Fatalf("remove db file: %v", err)
		}
	}
}

func doJSONRequest[T any](t *testing.T, client *http.Client, method, url string, payload any, result *T) int {
	t.Helper()

	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			t.Fatalf("encode payload: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}

	return resp.StatusCode
}

func TestBookstoreIntegration(t *testing.T) {
	server, _, cleanup := setupIntegrationServer(t)
	defer cleanup()

	client := server.Client()
	baseURL := server.URL

	var createdAuthor models.Author
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/authors", models.Author{Name: "Jane Doe"}, &createdAuthor); code != http.StatusCreated {
		t.Fatalf("expected 201 for author creation, got %d", code)
	}

	if createdAuthor.ID == 0 {
		t.Fatal("expected created author to have non-zero ID")
	}
	if createdAuthor.Category != "regular" {
		t.Fatalf("expected created author category regular, got %+v", createdAuthor)
	}

	var createdBook models.Book
	bookPayload := models.Book{
		Name:           "Example Book",
		PrintType:      "hardcoder",
		Paperback:      true,
		AuthorID:       createdAuthor.ID,
		PublisherHouse: "Acme Publishing",
		CopyAvailable:  true,
	}
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/books", bookPayload, &createdBook); code != http.StatusCreated {
		t.Fatalf("expected 201 for book creation, got %d", code)
	}

	var createdCustomer models.Customer
	customerPayload := models.Customer{
		Name:        "John Smith",
		Phone:       "123-456-7890",
		Email:       "john.smith@example.com",
		DateOfBirth: "1985-06-14",
	}
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/customers", customerPayload, &createdCustomer); code != http.StatusCreated {
		t.Fatalf("expected 201 for customer creation, got %d", code)
	}

	var createdPurchase models.Purchase
	purchasePayload := models.Purchase{
		CustomerID:    createdCustomer.ID,
		BookID:        createdBook.ID,
		PurchaseDate:  time.Now().UTC().Format(time.RFC3339),
		PurchasePrice: 29.99,
	}
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/purchases", purchasePayload, &createdPurchase); code != http.StatusCreated {
		t.Fatalf("expected 201 for purchase creation, got %d", code)
	}

	if createdPurchase.ID == 0 {
		t.Fatal("expected created purchase to have non-zero ID")
	}

	var fetchedAuthor models.Author
	if code := doJSONRequest(t, client, http.MethodGet, baseURL+"/authors/"+itoa(createdAuthor.ID), nil, &fetchedAuthor); code != http.StatusOK {
		t.Fatalf("expected 200 for author detail, got %d", code)
	}
	if fetchedAuthor.Category != "regular" {
		t.Fatalf("expected author detail category regular, got %+v", fetchedAuthor)
	}

	var authors []models.Author
	if code := doJSONRequest(t, client, http.MethodGet, baseURL+"/authors", nil, &authors); code != http.StatusOK {
		t.Fatalf("expected 200 for authors list, got %d", code)
	}
	if len(authors) != 1 || authors[0].Category != "regular" {
		t.Fatalf("unexpected authors list: %#v", authors)
	}

	var books []models.Book
	if code := doJSONRequest(t, client, http.MethodGet, baseURL+"/books", nil, &books); code != http.StatusOK {
		t.Fatalf("expected 200 for books list, got %d", code)
	}
	if len(books) != 1 || books[0].ID != createdBook.ID {
		t.Fatalf("unexpected books list: %#v", books)
	}

	var purchases []models.Purchase
	if code := doJSONRequest(t, client, http.MethodGet, baseURL+"/purchases", nil, &purchases); code != http.StatusOK {
		t.Fatalf("expected 200 for purchases list, got %d", code)
	}
	if len(purchases) != 1 || purchases[0].CustomerID != createdCustomer.ID {
		t.Fatalf("unexpected purchases list: %#v", purchases)
	}
}

func TestAuthorCategoryThresholdIntegration(t *testing.T) {
	server, _, cleanup := setupIntegrationServer(t)
	defer cleanup()

	client := server.Client()
	baseURL := server.URL

	var author models.Author
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/authors", models.Author{Name: "Premium Author"}, &author); code != http.StatusCreated {
		t.Fatalf("expected 201 for author creation, got %d", code)
	}

	var book models.Book
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/books", models.Book{
		Name:           "Hit Book",
		PrintType:      "hardcover",
		Paperback:      false,
		AuthorID:       author.ID,
		PublisherHouse: "Acme Publishing",
		CopyAvailable:  true,
	}, &book); code != http.StatusCreated {
		t.Fatalf("expected 201 for book creation, got %d", code)
	}

	var customer models.Customer
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/customers", models.Customer{Name: "Buyer"}, &customer); code != http.StatusCreated {
		t.Fatalf("expected 201 for customer creation, got %d", code)
	}

	purchasePayload := models.Purchase{
		CustomerID:    customer.ID,
		BookID:        book.ID,
		PurchaseDate:  time.Now().UTC().Format(time.RFC3339),
		PurchasePrice: 29.99,
	}

	for i := 0; i < 500; i++ {
		var purchase models.Purchase
		if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/purchases", purchasePayload, &purchase); code != http.StatusCreated {
			t.Fatalf("expected 201 for purchase %d, got %d", i+1, code)
		}
	}

	var fetchedAuthor models.Author
	if code := doJSONRequest(t, client, http.MethodGet, baseURL+"/authors/"+itoa(author.ID), nil, &fetchedAuthor); code != http.StatusOK {
		t.Fatalf("expected 200 for author detail at 500 sales, got %d", code)
	}
	if fetchedAuthor.Category != "regular" {
		t.Fatalf("expected regular category at 500 sales, got %+v", fetchedAuthor)
	}

	var purchase models.Purchase
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/purchases", purchasePayload, &purchase); code != http.StatusCreated {
		t.Fatalf("expected 201 for 501st purchase, got %d", code)
	}

	if code := doJSONRequest(t, client, http.MethodGet, baseURL+"/authors/"+itoa(author.ID), nil, &fetchedAuthor); code != http.StatusOK {
		t.Fatalf("expected 200 for author detail at 501 sales, got %d", code)
	}
	if fetchedAuthor.Category != "premium" {
		t.Fatalf("expected premium category at 501 sales, got %+v", fetchedAuthor)
	}
}

func TestAuthorCategoryRecalculationIntegration(t *testing.T) {
	server, database, cleanup := setupIntegrationServer(t)
	defer cleanup()

	client := server.Client()
	baseURL := server.URL

	var author models.Author
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/authors", models.Author{Name: "Recalc Author"}, &author); code != http.StatusCreated {
		t.Fatalf("expected 201 for author creation, got %d", code)
	}

	var book models.Book
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/books", models.Book{
		Name:           "Mutable Book",
		PrintType:      "paperback",
		Paperback:      true,
		AuthorID:       author.ID,
		PublisherHouse: "Acme Publishing",
		CopyAvailable:  true,
	}, &book); code != http.StatusCreated {
		t.Fatalf("expected 201 for book creation, got %d", code)
	}

	var customer models.Customer
	if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/customers", models.Customer{Name: "Buyer"}, &customer); code != http.StatusCreated {
		t.Fatalf("expected 201 for customer creation, got %d", code)
	}

	purchasePayload := models.Purchase{
		CustomerID:    customer.ID,
		BookID:        book.ID,
		PurchaseDate:  time.Now().UTC().Format(time.RFC3339),
		PurchasePrice: 29.99,
	}

	for i := 0; i < 501; i++ {
		var purchase models.Purchase
		if code := doJSONRequest(t, client, http.MethodPost, baseURL+"/purchases", purchasePayload, &purchase); code != http.StatusCreated {
			t.Fatalf("expected 201 for purchase %d, got %d", i+1, code)
		}
	}

	var fetchedAuthor models.Author
	if code := doJSONRequest(t, client, http.MethodGet, baseURL+"/authors/"+itoa(author.ID), nil, &fetchedAuthor); code != http.StatusOK {
		t.Fatalf("expected 200 for author detail before recalculation, got %d", code)
	}
	if fetchedAuthor.Category != "premium" {
		t.Fatalf("expected premium category before recalculation, got %+v", fetchedAuthor)
	}

	if _, err := database.Exec(
		"DELETE FROM customer_book_purchase WHERE id IN (SELECT id FROM customer_book_purchase WHERE book_id = ? ORDER BY id LIMIT 2)",
		book.ID,
	); err != nil {
		t.Fatalf("delete purchases: %v", err)
	}

	if code := doJSONRequest(t, client, http.MethodGet, baseURL+"/authors/"+itoa(author.ID), nil, &fetchedAuthor); code != http.StatusOK {
		t.Fatalf("expected 200 for author detail after recalculation, got %d", code)
	}
	if fetchedAuthor.Category != "regular" {
		t.Fatalf("expected regular category after recalculation, got %+v", fetchedAuthor)
	}
}

func itoa(id int64) string {
	return fmt.Sprintf("%d", id)
}
