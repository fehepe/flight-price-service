package server

import (
	"log"
	"os"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/internal/providers/amadeus"
	"github.com/fehepe/flight-price-service/internal/providers/mock"
	"github.com/fehepe/flight-price-service/internal/providers/serpapi"
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
	amadeusKey := mustEnv("AMADEUS_API_KEY")
	amadeusSecret := mustEnv("AMADEUS_API_SECRET")
	amadeusBaseURL := mustEnv("AMADEUS_API_BASE_URL")
	maxResults := mustEnv("MAX_FLIGHT_RESULTS_PER_CLIENT")

	// SerpAPI configuration
	serpapiKey := mustEnv("SER_API_KEY")
	serpapiBaseURL := mustEnv("SER_API_BASE_URL")

	return []providers.Provider{
		amadeus.New(amadeusKey, amadeusSecret, amadeusBaseURL, maxResults, nil),
		serpapi.New(serpapiKey, serpapiBaseURL, nil),
		mock.New(false),
	}
}
