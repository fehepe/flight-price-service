package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/internal/services/flight"
	"github.com/fehepe/flight-price-service/pkg/models"
	"github.com/fehepe/flight-price-service/pkg/utils"
)

// HealthCheck returns detailed service status.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s - HealthCheck", r.Method, r.URL.Path)

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"status":  "OK",
		"service": "flight-price-service",
		"version": "1.0.0",
	})
}

type FlightHandler struct {
	Providers []providers.Provider
}

func NewFlightHandler(providers []providers.Provider) *FlightHandler {
	return &FlightHandler{Providers: providers}
}

func (h *FlightHandler) GetFlights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search, err := extractFlightSearch(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var (
		allOffers []models.FlightOffer
		mu        sync.Mutex
		wg        sync.WaitGroup
	)

	for _, p := range h.Providers {
		wg.Add(1)
		go func(p providers.Provider) {
			defer wg.Done()
			offers, err := p.GetFlights(ctx, search)
			if err != nil {
				return
			}
			mu.Lock()
			allOffers = append(allOffers, offers...)
			mu.Unlock()
		}(p)
	}
	wg.Wait()

	if len(allOffers) == 0 {
		utils.RespondError(w, http.StatusNotFound, "no flight offers found")
		return
	}

	response := flight.BuildSearchResponse(allOffers)
	utils.RespondJSON(w, http.StatusOK, response)
}
