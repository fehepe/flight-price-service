package server

import (
	"log"
	"os"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/internal/providers/amadeus"
	"github.com/fehepe/flight-price-service/internal/providers/mock"
)

func MustLoadProviders() providers.Provider {
	apiKey := os.Getenv("AMADEUS_API_KEY")
	apiSecret := os.Getenv("AMADEUS_API_SECRET")
	baseURL := os.Getenv("AMADEUS_API_BASE_URL")
	maxResults := os.Getenv("MAX_FLIGHT_RESULTS_PER_CLIENT")

	if apiKey == "" || apiSecret == "" || baseURL == "" || maxResults == "" {
		log.Fatal("Missing one or more required Amadeus environment variables")
	}

	return providers.MultiProvider{
		Providers: []providers.Provider{
			amadeus.New(apiKey, apiSecret, baseURL, maxResults, nil),
			mock.New(false),
		},
	}
}
