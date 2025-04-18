package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv reads .env and panics on error.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment")
	}
}

// Get returns the value of the named env var or the fallback.
func Get(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
