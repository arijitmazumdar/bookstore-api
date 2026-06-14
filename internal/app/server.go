package app

import (
	"fmt"
	"net/http"

	"bookstore-api/internal/db"
)

func RunServer(addr, dbPath string) error {
	database, err := db.Open(dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if cerr := database.Close(); cerr != nil {
			fmt.Printf("warning: close database: %v\n", cerr)
		}
	}()

	if err := db.RunMigrations(database); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	router := NewRouter(database)
	return http.ListenAndServe(addr, router)
}
