package repository

import (
	"context"
	"database/sql"
	"time"

	 // replace with your actual module name
	 "github.com/raaj2493/Shortly/Backend/internals/models"
)

// UserRepository handles database operations for users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func ( r *UserRepository) CreateUser (ctx context.Context , user *models.User) error{
	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now().UTC()
}

