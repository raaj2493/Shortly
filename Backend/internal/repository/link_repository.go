package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/raaj2493/Shortly/Backend/internal/models"
)

// LinkRepository coordinates access mechanics to the URL storage vault
type LinkRepository struct {
	db *sql.DB
}

// NewLinkRepository configures a fresh operational repository terminal link
func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

// CreateLink saves an initial link record and returns a system auto-incremented serial sequence ID
func (r *LinkRepository) CreateLink(ctx context.Context, link *models.Link) (uint64, error) {
	query := `
		INSERT INTO links (original_url, user_id, click_count, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id;
	`
	var lastInsertID uint64
	err := r.db.QueryRowContext(ctx, query, link.OriginalURL, link.UserID, 0).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// UpdateShortCode patches the link record table with its matching calculated base62 code string token
func (r *LinkRepository) UpdateShortCode(ctx context.Context, id uint64, code string) error {
	query := `UPDATE links SET short_code = $1 WHERE id = $2;`
	_, err := r.db.ExecContext(ctx, query, code, id)
	return err
}

// GetLinkByCode scans our repository rows to pull out our original URL mapping destination target
func (r *LinkRepository) GetLinkByCode(ctx context.Context, code string) (*models.Link, error) {
	query := `
		SELECT id, short_code, original_url, user_id, click_count, created_at
		FROM links
		WHERE short_code = $1;
	`
	var link models.Link
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&link.ID, &link.ShortCode, &link.OriginalURL, &link.UserID, &link.ClickCount, &link.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("requested shortened link path does not exist in our system")
	} else if err != nil {
		return nil, err
	}

	return &link, nil
}

// IncrementClick records an isolated hit metrics log directly on our link row item counters
func (r *LinkRepository) IncrementClick(ctx context.Context, code string) error {
	query := `UPDATE links SET click_count = click_count + 1 WHERE short_code = $1;`
	_, err := r.db.ExecContext(ctx, query, code)
	return err
}