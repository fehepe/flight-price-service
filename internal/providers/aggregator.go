package providers

import (
	"context"
	"sync"

	"github.com/fehepe/flight-price-service/pkg/models"
)

// Provider is the interface that all flight data providers must implement.
type Provider interface {
	GetFlights(ctx context.Context, search models.FlightSearch) ([]models.FlightOffer, error)
}

// MultiProvider calls all providers concurrently and aggregates the results.
type MultiProvider struct {
	Providers []Provider
}

func (m MultiProvider) GetFlights(ctx context.Context, search models.FlightSearch) ([]models.FlightOffer, error) {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		all    []models.FlightOffer
		errors []error
	)

	for _, p := range m.Providers {
		wg.Add(1)
		go func(provider Provider) {
			defer wg.Done()
			flights, err := provider.GetFlights(ctx, search)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors = append(errors, err)
				return
			}
			all = append(all, flights...)
		}(p)
	}

	wg.Wait()

	return all, nil
}
