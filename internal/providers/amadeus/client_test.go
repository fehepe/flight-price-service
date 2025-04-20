package amadeus_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fehepe/flight-price-service/internal/providers/amadeus"
	"github.com/fehepe/flight-price-service/pkg/models"
)

func TestGetFlights_Success(t *testing.T) {
	t.Helper()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if strings.Contains(r.URL.Path, "/token") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"access_token":"mock-token","expires_in":3600}`))
			if err != nil {
				t.Fatalf("failed to write token response: %v", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"data": [{
				"itineraries": [{
					"duration": "PT3H30M",
					"segments": [{
						"departure": {"iataCode": "JFK", "at": "2025-05-02T10:00:00"},
						"arrival": {"iataCode": "LAX"}
					}]
				}],
				"price": {"total": "199.99"}
			}]
		}`))
		if err != nil {
			t.Fatalf("failed to write flight response: %v", err)
		}
	}))
	t.Cleanup(mockServer.Close)

	client := amadeus.New(
		"fake-api-key",
		"fake-api-secret",
		mockServer.URL,
		"10",
		mockServer.Client(),
	)

	search := models.FlightSearch{
		Origin:        "JFK",
		Destination:   "LAX",
		DepartureDate: time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC),
	}

	flights, err := client.GetFlights(context.Background(), search)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(flights) != 1 {
		t.Fatalf("expected 1 flight offer, got %d", len(flights))
	}

	if flights[0].Price != 199.99 {
		t.Errorf("expected price 199.99, got %f", flights[0].Price)
	}

	if flights[0].Duration != "PT3H30M" {
		t.Errorf("expected duration PT3H30M, got %s", flights[0].Duration)
	}
}

func TestGetFlights_TokenError(t *testing.T) {
	t.Helper()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`{"error":"invalid_client"}`))
		if err != nil {
			t.Fatalf("failed to write error response: %v", err)
		}
	}))
	t.Cleanup(mockServer.Close)

	client := amadeus.New(
		"invalid",
		"invalid",
		mockServer.URL,
		"10",
		mockServer.Client(),
	)

	search := models.FlightSearch{
		Origin:        "JFK",
		Destination:   "LAX",
		DepartureDate: time.Now().AddDate(0, 0, 1),
	}

	_, err := client.GetFlights(context.Background(), search)
	if err == nil || !strings.Contains(err.Error(), "auth failed") {
		t.Errorf("expected auth failed, got: %v", err)
	}
}
