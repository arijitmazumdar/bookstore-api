package models

type Book struct {
	ID             int64  `json:"id,omitempty"`
	Name           string `json:"name"`
	PrintType      string `json:"print_type"`
	Paperback      bool   `json:"paperback"`
	AuthorID       int64  `json:"author_id"`
	PublisherHouse string `json:"publisher_house"`
	CopyAvailable  bool   `json:"copy_available"`
}

type Author struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
}

type Customer struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}

type Purchase struct {
	ID            int64   `json:"id,omitempty"`
	CustomerID    int64   `json:"customer_id"`
	BookID        int64   `json:"book_id"`
	PurchaseDate  string  `json:"purchase_date"`
	PurchasePrice float64 `json:"purchase_price"`
}
