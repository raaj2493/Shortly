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

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, now, now).Scan(&user.ID)
	if err != nil {
		return err
	}
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *UserRepository) GetUserByEmail (ctx context.Context , email string ) (*models.User , error){
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil , nil
		}
		return nil , err
	}
	return &user , nil
}

