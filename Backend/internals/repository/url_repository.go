package repository

import (
	"database/sql"
	"context"

	"github.com/raaj2493/Shortly/Backend/internals/models"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) CreateURL(ctx context.Context, originalURL , shortCode string) (*models.URL, error) {
	query := `
		INSERT INTO urls (original_url, short_code)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	var url models.URL
	err := r.db.QueryRowContext(ctx, query, originalURL, shortCode).Scan(&url.ID, &url.CreatedAt)
	if err != nil {
		return nil, err
	}
	url.OriginalURL = originalURL
	url.ShortCode = shortCode
	return &url, nil
}

func (r *URLRepository) GetURL(ctx context.Context, shortCode string) (*models.URL, error) {
	query := `
		SELECT id, original_url, short_code, created_at
		FROM urls
		WHERE short_code = $1
	`
	var url models.URL
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &url, nil
}