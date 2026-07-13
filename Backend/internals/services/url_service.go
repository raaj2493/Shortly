package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/raaj2493/Shortly/Backend/internals/repository"
	"github.com/redis/go-redis/v9"
)

const (
	shortCodeLength = 6
	cachePrefix     = "url:"
	cacheTTL        = 1 * time.Hour
	maxRetries      = 5
)

type URLService struct {
	repo  *repository.URLRepository
	redis *redis.Client
}

func NewURLService(repo *repository.URLRepository, redis *redis.Client) *URLService {
	return &URLService{repo: repo, redis: redis}
}

// Shorten generates a unique short code and stores the mapping.
func (s *URLService) Shorten(ctx context.Context, originalURL string) (string, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		code, err := generateShortCode()
		if err != nil {
			return "", fmt.Errorf("failed to generate short code: %w", err)
		}

		url, err := s.repo.CreateURL(ctx, originalURL, code)
		if err == nil {
			// Optionally cache the new mapping immediately
		 s.redis.Set(ctx, cachePrefix+code, originalURL, cacheTTL)
			return url.ShortCode, nil
		}

		// Check if it's a unique violation (PostgreSQL error code 23505)
		if isDuplicateKeyError(err) {
			continue // try another code
		}
		return "", fmt.Errorf("failed to save URL: %w", err)
	}
	return "", errors.New("failed to generate a unique short code after multiple attempts")
}


// GetOriginalURL retrieves the original URL, using Redis cache first.
func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	// 1. Check Redis cache
	cacheKey := cachePrefix + shortCode
	originalURL, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		return originalURL, nil
	}
	if !errors.Is(err, redis.Nil) {
		// Real error (not just cache miss)
		return "", fmt.Errorf("redis error: %w", err)
	}

	// 2. Cache miss – query PostgreSQL
	url, err := s.repo.GetURL(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}
	if url == nil {
		return "", errors.New("URL not found")
	}

	// 3. Store in Redis for future requests
	// We ignore error on set; caching is an optimisation, not critical
	_ = s.redis.Set(ctx, cacheKey, url.OriginalURL, cacheTTL).Err()

	return url.OriginalURL, nil
}




func generateShortCode() (string , error){
	bytes := make([]byte, shortCodeLength)
	if _, err:= rand.Read(bytes)
	err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes)[:shortCodeLength], nil
}

func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}