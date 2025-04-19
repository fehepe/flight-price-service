package flight

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/pkg/models"
)

// FetchAllFlightOffers retrieves and merges flight offers from all providers concurrently.
func FetchAllFlightOffers(ctx context.Context, providerList []providers.Provider, search models.FlightSearch) ([]models.FlightOffer, error) {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		all  []models.FlightOffer
		errs []error
	)

	for _, p := range providerList {
		wg.Add(1)
		go func(pr providers.Provider) {
			defer wg.Done()
			offers, err := pr.GetFlights(ctx, search)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			all = append(all, offers...)
		}(p)
	}

	wg.Wait()

	var combinedErr error
	if len(errs) > 0 {
		errMsgs := make([]string, len(errs))
		for i, e := range errs {
			errMsgs[i] = e.Error()
		}
		combinedErr = fmt.Errorf("provider errors: %s", strings.Join(errMsgs, "; "))
	}

	return all, combinedErr
}
