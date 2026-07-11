package config

import (
    "os"
    "strconv"
)

type Config struct {
    ServerPort   string
    DatabaseDSN  string
    RedisAddr    string
    RedisPass    string
    JWTSecret    string
}

func Load() *Config{
	return &Config{
        ServerPort:  getEnv("SERVER_PORT", "8080"),
        DatabaseDSN: getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"),
        RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
        RedisPass:   getEnv("REDIS_PASSWORD", ""),
        JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
    }
}

func getEnv(key , fallback string) string  {
	if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

// Helper to convert string to int (not used now, but handy for future)
func getEnvInt(key string, fallback int) int {
    if value, ok := os.LookupEnv(key); ok {
        if i, err := strconv.Atoi(value); err == nil {
            return i
        }
    }
    return fallback
}