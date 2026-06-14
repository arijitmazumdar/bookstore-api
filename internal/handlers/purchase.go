package handlers

import (
	"encoding/json"
	"net/http"

	"bookstore-api/internal/models"
)

func (h *Handler) PurchasesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listPurchases(w, r)
	case http.MethodPost:
		h.createPurchase(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listPurchases(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, customer_id, book_id, purchase_date, purchase_price FROM customer_book_purchase")
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	purchases := []models.Purchase{}
	for rows.Next() {
		var purchase models.Purchase
		if err := rows.Scan(&purchase.ID, &purchase.CustomerID, &purchase.BookID, &purchase.PurchaseDate, &purchase.PurchasePrice); err != nil {
			h.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		purchases = append(purchases, purchase)
	}

	h.writeJSON(w, http.StatusOK, purchases)
}

func (h *Handler) createPurchase(w http.ResponseWriter, r *http.Request) {
	var purchase models.Purchase
	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	result, err := h.db.Exec("INSERT INTO customer_book_purchase (customer_id, book_id, purchase_date, purchase_price) VALUES (?, ?, ?, ?)", purchase.CustomerID, purchase.BookID, purchase.PurchaseDate, purchase.PurchasePrice)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	purchase.ID, _ = result.LastInsertId()
	h.writeJSON(w, http.StatusCreated, purchase)
}
