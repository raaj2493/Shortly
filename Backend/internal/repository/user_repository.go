package repository

import (
	"database/sql"
	"context"
	"time"

	"github.com/raaj2493/Shortly/Backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository (db *sql.DB) *UserRepository{
	return &UserRepository{db: db}
}

// CreateUser inserts a brand new user into the database
func(r *UserRepository) CreateUser(ctx context.Context , user *models.User) error {
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
func(r *UserRepository)GetUserByEmail(ctx context.Context, email string ) (*models.User , error ){
	ctx , cancel := context.WithTimeout(ctx , 3*time.Second)
	defer cancel()

	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx , query , email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err // Will return sql.ErrNoRows if user isn't found
	}

	return &user, nil
}