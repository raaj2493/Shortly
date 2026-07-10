package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application.
type Config struct {
	// Server
	ServerPort int

	// PostgreSQL
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string // e.g., "disable" for local dev

	// Redis (we'll use later)
	RedisHost string
	RedisPort int

	// Logger
	LogLevel  string // "debug", "info", "warn", "error"
	LogFormat string // "json" or "text"
}


func Load()(*Config , error){
	getEnv := func(key , defaultVal string) string {
		if val := os.Getenv(key);
		val != ""{
			return val
		}
		return defaultVal
	}
	getEnvInt := func (key , defaultVal int) (int , error){
		val := os.Getenv(key)
		if val == "" {
			return defaultVal, nil
		}

		intVal, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("invalid integer for %s: %w", key, err)
		}
		return intVal, nil
	}
	serverPort, err := getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	dbPort, err := getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	redisPort, err := getEnvInt("REDIS_PORT", 6379)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerPort: serverPort,
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "url_shortener"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  redisPort,
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		LogFormat:  getEnv("LOG_FORMAT", "json"),
	}

	// Basic validation (optional but helpful)
	if cfg.DBUser == "" {
		return nil, errors.New("DB_USER is required")
	}

	return cfg, nil

}