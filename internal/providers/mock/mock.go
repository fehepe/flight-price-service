package mock

import (
	"context"
	"errors"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/pkg/models"
)

// MockProvider returns dummy flight data for testing and local dev.
type MockProvider struct {
	ShouldFail bool
}

func New(shouldFail bool) providers.Provider {
	return &MockProvider{ShouldFail: shouldFail}
}

// GetFlights returns mock flights or simulates an error.
func (m *MockProvider) GetFlights(ctx context.Context, search models.FlightSearch) ([]models.FlightOffer, error) {
	if m.ShouldFail {
		return nil, errors.New("mock provider error")
	}

	date := search.DepartureDate.Format("2006-01-02")

	offers := []models.FlightOffer{
		{
			Provider:    "MockAir",
			Price:       80.00,
			Duration:    "PT16H30M",
			Origin:      search.Origin,
			Destination: search.Destination,
			Date:        date,
		},
		{
			Provider:    "MockExpress",
			Price:       150.00,
			Duration:    "PT3H0M",
			Origin:      search.Origin,
			Destination: search.Destination,
			Date:        date,
		},
	}

	return offers, nil
}
