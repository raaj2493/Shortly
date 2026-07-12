package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/raaj2493/Shortly/Backend/internals/models"
	"github.com/raaj2493/Shortly/Backend/internals/repository"

	"github.com/redis/go-redis/v9"
)

// URLService handles URL shortening logic.
type URLService struct {
	urlRepo   *repository.URLRepository
	redis     *redis.Client
	ctx       context.Context
	codeLen   int
	maxRetries int
}

// NewURLService creates a new instance.
func NewURLService(urlRepo *repository.URLRepository, redis *redis.Client) *URLService {
	return &URLService{
		urlRepo:    urlRepo,
		redis:      redis,
		ctx:        context.Background(),
		codeLen:    6,
		maxRetries: 3,
	}
}

// generateShortCode creates a random 6-character string using crypto/rand.
func (s *URLService) generateShortCode() (string, error) {
	bytes := make([]byte, s.codeLen)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// Use base64 URL encoding without padding; it gives alphanumeric + '-' and '_'
	// but we want only alphanumeric to keep it simple. We'll map to a custom alphabet.
	// Safer: use a fixed alphabet and random bytes modulo.
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, s.codeLen)
	for i := 0; i < s.codeLen; i++ {
		// Use a random byte modulo the alphabet size; this introduces a tiny bias but is fine.
		// For production, use crypto/rand with rejection sampling, but this is okay for learning.
		idx := int(bytes[i]) % len(alphabet)
		code[i] = alphabet[idx]
	}
	return string(code), nil
}

// CreateShortURL creates a short URL for a given original URL and user ID.
// It generates a unique short code, saves to DB, and caches in Redis.
func (s *URLService) CreateShortURL(originalURL string, userID int64) (*models.URL, error) {
	// Try to generate a unique short code with retries.
	var shortCode string
	var err error
	for i := 0; i < s.maxRetries; i++ {
		shortCode, err = s.generateShortCode()
		if err != nil {
			return nil, fmt.Errorf("generate short code: %w", err)
		}
		exists, err := s.urlRepo.ShortCodeExists(shortCode)
		if err != nil {
			return nil, fmt.Errorf("check existence: %w", err)
		}
		if !exists {
			break // unique code found
		}
		if i == s.maxRetries-1 {
			return nil, errors.New("failed to generate unique short code after retries")
		}
	}

	url := &models.URL{
		UserID:      userID,
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now().UTC(),
	}

	// Save to database
	if err := s.urlRepo.CreateURL(url); err != nil {
		return nil, fmt.Errorf("save URL: %w", err)
	}

	// Cache it in Redis with a TTL (e.g., 1 hour)
	cacheKey := fmt.Sprintf("url:%s", shortCode)
	if err := s.redis.Set(s.ctx, cacheKey, originalURL, time.Hour).Err(); err != nil {
		// Log the error but don't fail the whole operation; caching is optional.
		// In a real app, you'd log this.
	}

	return url, nil
}

// GetOriginalURL retrieves the original URL from cache or database.
func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	// 1. Try Redis cache.
	cacheKey := fmt.Sprintf("url:%s", shortCode)
	original, err := s.redis.Get(s.ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		return original, nil
	}
	if !errors.Is(err, redis.Nil) {
		// Some other Redis error – we could log it, but we'll fall back to DB.
		// For learning, we just ignore and proceed.
	}

	// 2. Cache miss: query PostgreSQL.
	url, err := s.urlRepo.GetURLByShortCode(shortCode)
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}
	if url == nil {
		return "", errors.New("short URL not found")
	}

	// 3. (Optional) Update cache for future requests.
	// We'll set the cache with TTL.
	if err := s.redis.Set(s.ctx, cacheKey, url.OriginalURL, time.Hour).Err(); err != nil {
		// Log but don't fail.
	}

	return url.OriginalURL, nil
}