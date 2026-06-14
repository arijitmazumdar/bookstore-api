package main

import (
	"log"
	"os"

	"bookstore-api/internal/app"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	dbPath := "bookstore.db"
	if env := os.Getenv("DATABASE_PATH"); env != "" {
		dbPath = env
	}

	log.Printf("starting bookstore API on %s using database %s", addr, dbPath)
	if err := app.RunServer(addr, dbPath); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
