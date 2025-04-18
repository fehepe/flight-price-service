package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fehepe/flight-price-service/pkg/models"
)

// respondJSON sets headers and writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("[respondJSON] error encoding response: %v", err)
	}
}

// HealthCheck returns a simple OK status to indicate the service is healthy.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s - HealthCheck", r.Method, r.URL.Path)
	respondJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

// GetFlights returns placeholder flight offers based on a search.
func GetFlights(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s - GetFlights", r.Method, r.URL.Path)

	// TODO: Replace stub data with real search logic or service integration
	offers := []models.FlightOffer{
		{
			Provider:    "Airline A",
			Price:       100.0,
			Duration:    2 * time.Hour,
			Origin:      "NYC",
			Destination: "LAX",
			Date:        "2023-10-01",
		},
		{
			Provider:    "Airline A",
			Price:       120.0,
			Duration:    3 * time.Hour,
			Origin:      "NYC",
			Destination: "LAX",
			Date:        "2023-10-02",
		},
		{
			Provider:    "Airline B",
			Price:       150.0,
			Duration:    1 * time.Hour,
			Origin:      "NYC",
			Destination: "LAX",
			Date:        "2023-10-01",
		},
		{
			Provider:    "Airline B",
			Price:       130.0,
			Duration:    2 * time.Hour,
			Origin:      "NYC",
			Destination: "LAX",
			Date:        "2023-10-02",
		},
	}

	response := models.SearchResponse{
		Providers: map[string][]models.FlightOffer{
			"Airline A": offers[:2],
			"Airline B": offers[2:],
		},
		Cheapest: models.FlightOffer{},
		Fastest:  models.FlightOffer{},
	}

	// Find the cheapest and fastest
	for _, offer := range offers {
		if response.Cheapest.Price == 0 || offer.Price < response.Cheapest.Price {
			response.Cheapest = offer
		}
		if response.Fastest.Duration == 0 || offer.Duration < response.Fastest.Duration {
			response.Fastest = offer
		}
	}

	respondJSON(w, http.StatusOK, response)
}
