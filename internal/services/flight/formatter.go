package flight

import (
	"strconv"
	"strings"
	"time"

	"github.com/fehepe/flight-price-service/pkg/models"
)

func BuildSearchResponse(offers []models.FlightOffer) models.SearchResponse {
	var (
		cheapest    = offers[0]
		fastest     = offers[0]
		minDuration = mustParseISODuration(offers[0].Duration)
	)

	providerMap := make(map[string][]models.FlightOffer)
	for _, offer := range offers {
		if offer.Price < cheapest.Price {
			cheapest = offer
		}

		if dur := mustParseISODuration(offer.Duration); dur < minDuration {
			fastest = offer
			minDuration = dur
		}

		providerMap[offer.Provider] = append(providerMap[offer.Provider], offer)
	}

	return models.SearchResponse{
		Cheapest:  cheapest,
		Fastest:   fastest,
		Providers: providerMap,
	}
}

func mustParseISODuration(iso string) time.Duration {
	iso = strings.TrimPrefix(strings.ToUpper(iso), "PT")

	var dur time.Duration
	if hIdx := strings.Index(iso, "H"); hIdx != -1 {
		hours, err := strconv.Atoi(iso[:hIdx])
		if err != nil {
			return 0
		}
		dur += time.Duration(hours) * time.Hour
		iso = iso[hIdx+1:]
	}
	if mIdx := strings.Index(iso, "M"); mIdx != -1 {
		mins, err := strconv.Atoi(iso[:mIdx])
		if err != nil {
			return 0
		}
		dur += time.Duration(mins) * time.Minute
	}
	return dur
}
