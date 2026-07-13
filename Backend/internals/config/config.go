package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	RedisURL    string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://shortener:shortener_pass@localhost:5432/url_shortener?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis:6379"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}