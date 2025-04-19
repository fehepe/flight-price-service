package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fehepe/flight-price-service/internal/cache"
	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/internal/services/flight"
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
	providers []providers.Provider
	cache     cache.FlightCacher
}

func NewFlightHandler(providers []providers.Provider, cache cache.FlightCacher) *FlightHandler {
	return &FlightHandler{providers: providers, cache: cache}
}

func (h *FlightHandler) GetFlights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search, err := extractFlightSearch(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := fmt.Sprintf("%s:%s:%s:%d:%t",
		search.Origin,
		search.Destination,
		search.DepartureDate.Format("2006-01-02"),
		search.Adults,
		search.NonStop,
	)

	cachedOffers, found, err := h.cache.Get(ctx, cacheKey)
	if err != nil {
		log.Printf("%s %s cache set get: %v\n", r.Method, r.RequestURI, err)
		utils.RespondError(w, http.StatusInternalServerError, "cache error")
		return
	}
	if found {
		utils.RespondJSON(w, http.StatusOK, flight.BuildSearchResponse(cachedOffers))
		return
	}

	offers := flight.FetchAllFlightOffers(ctx, h.providers, search)
	if err := h.cache.Set(ctx, cacheKey, offers); err != nil {
		log.Printf("%s %s cache set error: %v\n", r.Method, r.RequestURI, err)
	}

	utils.RespondJSON(w, http.StatusOK, flight.BuildSearchResponse(offers))
}
