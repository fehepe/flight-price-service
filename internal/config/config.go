package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
}

// Get returns the value of the named env var or the fallback.
func Get(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

// getEnvInt reads an environment variable into an integer or returns the fallback value.
func GetEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if iv, err := strconv.Atoi(v); err == nil {
			return iv
		}
	}
	return fallback
}
