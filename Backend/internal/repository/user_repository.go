package repository

import (
	"database/sql"
	"context"
	"time"

	"github.com/raaj2493/Shortly/Backend/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo (db *sql.DB) *UserRepo{
	return &UserRepo{db: db}
}

// CreateUser inserts a brand new user into the database
func(r *UserRepo) CreateUser(ctx context.Context , user *models.User) error {
     // Set a quick safety timer. If the database takes longer than 3 seconds, cancel it.
	 ctx , cancel := context.WithTimeout(ctx , 3*time.Second)
	 defer cancel()

	 query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx , query , user.Username, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

		return err 
}

// GetUserByEmail searches for a user by their unique email
func(r *UserRepo)GetUserByEmail(ctx context.Context, email string ) (*models.User , error ){
	ctx , cancel := context.WithTimeout(ctx , 3*time.Second)
	defer cancel()

	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	var user models.User
	err := r.db.
}