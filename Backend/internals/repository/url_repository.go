package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/raaj2493/Shortly/Backend/internals/models"
)

// URLRepository handles URL storage.
type URLRepository struct {
	db *sql.DB
}

// NewURLRepository creates a new instance.
func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

// CreateURL inserts a new URL mapping into the database.
func (r *URLRepository) CreateURL(url *models.URL) error {
	query := `
		INSERT INTO urls (user_id, original_url, short_code, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		url.UserID,
		url.OriginalURL,
		url.ShortCode,
		url.CreatedAt,
	).Scan(&url.ID)
	if err != nil {
		return fmt.Errorf("CreateURL: %w", err)
	}
	return nil
}

// GetURLByShortCode retrieves a URL by its short code.
func (r *URLRepository) GetURLByShortCode(shortCode string) (*models.URL, error) {
	query := `
		SELECT id, user_id, original_url, short_code, created_at
		FROM urls
		WHERE short_code = $1
	`
	var url models.URL
	err := r.db.QueryRow(query, shortCode).Scan(
		&url.ID,
		&url.UserID,
		&url.OriginalURL,
		&url.ShortCode,
		&url.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found, but not an error
		}
		return nil, fmt.Errorf("GetURLByShortCode: %w", err)
	}
	return &url, nil
}

// ShortCodeExists checks if a given short code is already taken.
func (r *URLRepository) ShortCodeExists(shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	var exists bool
	err := r.db.QueryRow(query, shortCode).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ShortCodeExists: %w", err)
	}
	return exists, nil
}