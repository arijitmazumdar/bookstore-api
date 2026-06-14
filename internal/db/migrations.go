package db

import (
	"database/sql"
)

func RunMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS authors (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS customers (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name TEXT NOT NULL,
		    phone TEXT,
		    email TEXT,
		    date_of_birth TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS books (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name TEXT NOT NULL,
		    print_type TEXT NOT NULL,
		    paperback BOOLEAN NOT NULL CHECK (paperback IN (0,1)),
		    author_id INTEGER NOT NULL,
		    publisher_house TEXT,
		    copy_available BOOLEAN NOT NULL CHECK (copy_available IN (0,1)),
		    FOREIGN KEY(author_id) REFERENCES authors(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS customer_book_purchase (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    customer_id INTEGER NOT NULL,
		    book_id INTEGER NOT NULL,
		    purchase_date TEXT NOT NULL,
		    purchase_price REAL NOT NULL,
		    FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE,
		    FOREIGN KEY(book_id) REFERENCES books(id) ON DELETE CASCADE
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}
