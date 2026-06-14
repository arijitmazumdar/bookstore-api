package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=1")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		if cerr := db.Close(); cerr != nil {
			return nil, fmt.Errorf("ping error: %v; close error: %w", err, cerr)
		}
		return nil, err
	}

	return db, nil
}
