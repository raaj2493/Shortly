package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	
)

type Config struct {
	Port       string
	Env        string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBHost     string
	JWTSecret  string
}



func LoadConfig()*Config{

	err := godotenv.Load()
	if err != nil {
		log.Println("Note: No .env file found, using system environment variables instead.")
	}
	
	return &Config{
		Port:       getEnv("PORT", "8080"),
		Env:        getEnv("ENV", "development"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "url_shortener"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		JWTSecret: getEnv("JWT_SECRET", "fallback_secret"),
	}
}

// Helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == "" {
			log.Fatalf("Critical environment variable missing: %s", key)
		}
		return defaultValue
	}
	return value
}