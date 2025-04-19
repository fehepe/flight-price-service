package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/internal/providers/mock"
	"github.com/fehepe/flight-price-service/pkg/models"
)

func TestGetFlights(t *testing.T) {
	h := NewFlightHandler([]providers.Provider{
		&mock.MockProvider{},
	})
	today := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "valid request",
			query:      "/flights/search?origin=JFK&destination=LAX&date=" + today,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing origin",
			query:      "/flights/search?destination=LAX&date=" + today,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid date format",
			query:      "/flights/search?origin=JFK&destination=LAX&date=02-01-2025",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "past departure date",
			query:      "/flights/search?origin=JFK&destination=LAX&date=2000-01-01",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid IATA code",
			query:      "/flights/search?origin=JF&destination=LAX&date=" + today,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid non_stop",
			query:      "/flights/search?origin=JFK&destination=LAX&date=" + today + "&non_stop=yes",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid adults",
			query:      "/flights/search?origin=JFK&destination=LAX&date=" + today + "&adults=ten",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "too many adults",
			query:      "/flights/search?origin=JFK&destination=LAX&date=" + today + "&adults=10",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			rec := httptest.NewRecorder()
			h.GetFlights(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.wantStatus, res.StatusCode)
			}

			if res.StatusCode == http.StatusOK {
				var resp models.SearchResponse
				if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
					t.Errorf("%s: failed to decode response: %v", tt.name, err)
				}
				if resp.Cheapest.Price != 80.0 {
					t.Errorf("%s: expected cheapest price 80.0, got %.2f", tt.name, resp.Cheapest.Price)
				}
				if !strings.EqualFold(resp.Fastest.Duration, "PT3H0M") {
					t.Errorf("%s: expected fastest duration PT3H0M, got %s", tt.name, resp.Fastest.Duration)
				}
			}
		})
	}
}
