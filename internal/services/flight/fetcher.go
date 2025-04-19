package flight

import (
	"context"
	"log"
	"sync"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/pkg/models"
)

// FetchAllFlightOffers retrieves and merges flight offers from all providers concurrently.
func FetchAllFlightOffers(ctx context.Context, providerList []providers.Provider, search models.FlightSearch) []models.FlightOffer {
	var (
		allOffers []models.FlightOffer
		mu        sync.Mutex
		wg        sync.WaitGroup
	)

	for _, provider := range providerList {
		wg.Add(1)
		go func(p providers.Provider) {
			defer wg.Done()
			offers, err := p.GetFlights(ctx, search)
			if err != nil {
				log.Printf("provider error: %v", err)
				return
			}
			mu.Lock()
			allOffers = append(allOffers, offers...)
			mu.Unlock()
		}(provider)
	}

	wg.Wait()
	return allOffers
}
