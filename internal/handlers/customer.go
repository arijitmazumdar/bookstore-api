package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bookstore-api/internal/models"
)

func (h *Handler) CustomersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listCustomers(w, r)
	case http.MethodPost:
		h.createCustomer(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) CustomerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getCustomer(w, r, id)
	case http.MethodPut:
		h.updateCustomer(w, r, id)
	case http.MethodDelete:
		h.deleteCustomer(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, phone, email, date_of_birth FROM customers")
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	customers := []models.Customer{}
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Email, &customer.DateOfBirth); err != nil {
			h.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		customers = append(customers, customer)
	}

	h.writeJSON(w, http.StatusOK, customers)
}

func (h *Handler) getCustomer(w http.ResponseWriter, r *http.Request, id int64) {
	var customer models.Customer
	row := h.db.QueryRow("SELECT id, name, phone, email, date_of_birth FROM customers WHERE id = ?", id)
	if err := row.Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Email, &customer.DateOfBirth); err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "customer not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, customer)
}

func (h *Handler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	result, err := h.db.Exec("INSERT INTO customers (name, phone, email, date_of_birth) VALUES (?, ?, ?, ?)", customer.Name, customer.Phone, customer.Email, customer.DateOfBirth)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	customer.ID, _ = result.LastInsertId()
	h.writeJSON(w, http.StatusCreated, customer)
}

func (h *Handler) updateCustomer(w http.ResponseWriter, r *http.Request, id int64) {
	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	customer.ID = id
	_, err := h.db.Exec("UPDATE customers SET name = ?, phone = ?, email = ?, date_of_birth = ? WHERE id = ?", customer.Name, customer.Phone, customer.Email, customer.DateOfBirth, id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, customer)
}

func (h *Handler) deleteCustomer(w http.ResponseWriter, r *http.Request, id int64) {
	_, err := h.db.Exec("DELETE FROM customers WHERE id = ?", id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
