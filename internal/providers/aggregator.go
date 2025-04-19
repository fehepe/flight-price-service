package providers

import (
	"context"

	"github.com/fehepe/flight-price-service/pkg/models"
)

// Provider is the interface that all flight data providers must implement.
type Provider interface {
	GetFlights(ctx context.Context, search models.FlightSearch) ([]models.FlightOffer, error)
}
