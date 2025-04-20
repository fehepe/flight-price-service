package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/fehepe/flight-price-service/pkg/models"
)

var iataRegex = regexp.MustCompile(`^[A-Z]{3}$`)

func extractFlightSearch(r *http.Request) (models.FlightSearch, error) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")
	dateStr := r.URL.Query().Get("date")

	if origin == "" || destination == "" || dateStr == "" {
		return models.FlightSearch{}, errors.New("missing required query parameters: origin, destination, date")
	}

	if !iataRegex.MatchString(origin) || !iataRegex.MatchString(destination) {
		return models.FlightSearch{}, errors.New("invalid IATA code format; expected 3 uppercase letters")
	}

	departureDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return models.FlightSearch{}, errors.New("invalid date format; expected YYYY-MM-DD")
	}

	if departureDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return models.FlightSearch{}, errors.New("date cannot be in the past")
	}

	return models.FlightSearch{
		Origin:        origin,
		Destination:   destination,
		DepartureDate: departureDate,
	}, nil
}
