package database


import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Explicitly loads the Postgres driver
	"github.com/raaj2493/Shortly/Backend/internal/config"
)

// InitDB sets up the PostgreSQL connection pool
func InitDB(cfg *config.Config) *sql.DB{
	// 1. Build the Data Source Name (DSN) connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// 2. Open a connection to the database
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection string: %v", err)
	}

	// 3. Set connection pool rules (Crucial for production stability)
	db.SetMaxOpenConns(25)                 // Max simultaneous connections active
	db.SetMaxIdleConns(25)                 // Max idle connections kept alive in the pool
	db.SetConnMaxLifetime(5 * time.Minute) // Recycles connections older than 5 mins


	// 4. Ping the database to make sure the credentials and network actually work
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	log.Println("✅ Database connection pool successfully initialized!")
	return db
}