package server

import (
	"log"
	"os"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/internal/providers/amadeus"
	"github.com/fehepe/flight-price-service/internal/providers/mock"
	"github.com/fehepe/flight-price-service/internal/providers/serpapi"
	"github.com/fehepe/flight-price-service/internal/secret"
)

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return value
}

func MustLoadProviders() []providers.Provider {
	// Amadeus configuration
	amadeusBaseURL := mustEnv("AMADEUS_API_BASE_URL")
	maxResults := mustEnv("MAX_FLIGHT_RESULTS_PER_CLIENT")

	// SerpAPI configuration
	serApiBaseURL := mustEnv("SER_API_BASE_URL")

	// Load encrypted credentials.json
	creds, err := secret.LoadCreds("credentials.json")
	if err != nil {
		log.Fatalf("cannot load credentials: %v", err)
	}

	// Validate loaded credentials
	if creds.AmadeusAPIKey == "" || creds.AmadeusAPISecret == "" {
		log.Fatal("amadeus credentials (API key & secret) must not be empty")
	}
	if creds.SerAPIKey == "" {
		log.Fatal("SerpAPI credential (API key) must not be empty")
	}

	return []providers.Provider{
		amadeus.New(creds.AmadeusAPIKey, creds.AmadeusAPISecret, amadeusBaseURL, maxResults, nil),
		serpapi.New(creds.SerAPIKey, serApiBaseURL, nil),
		mock.New(false),
	}
}
