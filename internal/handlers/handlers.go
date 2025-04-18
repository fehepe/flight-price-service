package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/fehepe/flight-price-service/pkg/models"
	"github.com/fehepe/flight-price-service/pkg/utils"
)

// HealthCheck returns a simple status OK.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s - HealthCheck", r.Method, r.URL.Path)
	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

// GetFlights returns placeholder flight offers.
func GetFlights(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s - GetFlights", r.Method, r.URL.Path)

	// TODO: Replace stub data with real provider integration.
	offers := []models.FlightOffer{
		{Provider: "Airline A", Price: 100.0, Duration: 2 * time.Hour, Origin: "NYC", Destination: "LAX", Date: "2023-10-01"},
		{Provider: "Airline B", Price: 150.0, Duration: 1 * time.Hour, Origin: "NYC", Destination: "LAX", Date: "2023-10-01"},
	}

	// Compute cheapest and fastest
	resp := models.SearchResponse{Providers: map[string][]models.FlightOffer{"default": offers}}
	for _, o := range offers {
		if resp.Cheapest.Price == 0 || o.Price < resp.Cheapest.Price {
			resp.Cheapest = o
		}
		if resp.Fastest.Duration == 0 || o.Duration < resp.Fastest.Duration {
			resp.Fastest = o
		}
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}
