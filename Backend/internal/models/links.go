package models

import "time"

// Link defines our primary database record parameters blueprint for shortened URLs
type Link struct {
	ID          uint64    `json:"id" db:"id"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	UserID      string    `json:"user_id" db:"user_id"`
	ClickCount  int64     `json:"click_count" db:"click_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}