package models

import "time"

// URL represents a shortened URL mapping.
type URL struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
}

