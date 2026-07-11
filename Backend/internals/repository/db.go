package repository


import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq" // PostgreSQL driver
)


var DB *sql.DB

func InitDB(dsn string){
	var err error
	DB , err = sql.Open("postgres" , dsn)

	if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }

	 // Connection pool settings
    DB.SetMaxOpenConns(25)
    DB.SetMaxIdleConns(5)

	if err = DB.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    fmt.Println("Database connection established")

}