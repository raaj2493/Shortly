package repository

import (
	"database/sql"
	"log"
"fmt"
	_ "github.com/lib/pq"
)


func NewDB(databaseURL string) *sql.DB {
	db , err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	return db
}