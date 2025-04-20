package serpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fehepe/flight-price-service/pkg/models"
)

func TestGetFlights_Success(t *testing.T) {
	now := time.Date(2025, 4, 21, 0, 0, 0, 0, time.UTC)
	sel := models.FlightSearch{
		Origin:        "AAA",
		Destination:   "BBB",
		DepartureDate: now,
	}

	handler := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.SerAPIResponse{
			BestFlights: []models.FlightOption{{
				Flights: []models.FlightSegment{{
					DepartureAirport: models.AirportInfo{ID: "AAA", Time: "2025-04-21 10:00"},
					ArrivalAirport:   models.AirportInfo{ID: "CCC", Time: "2025-04-21 14:30"},
				}},
				TotalDuration: 270,
				Price:         100,
			}},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		data, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		if _, err := w.Write(data); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer handler.Close()

	client := New("key", handler.URL, handler.Client())
	offers, err := client.GetFlights(context.Background(), sel)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(offers) != 1 {
		t.Fatalf("expected 1 offer, got %d", len(offers))
	}

	o := offers[0]
	if o.Provider != providerName {
		t.Errorf("expected provider %q, got %q", providerName, o.Provider)
	}
	if o.Price != 100 {
		t.Errorf("expected price 100, got %v", o.Price)
	}
	if o.Duration != "PT4H30M" {
		t.Errorf("expected duration PT4H30M, got %q", o.Duration)
	}
	if o.Origin != "AAA" || o.Destination != "CCC" {
		t.Errorf("expected route AAA->CCC, got %s->%s", o.Origin, o.Destination)
	}
	if o.Date != "2025-04-21" {
		t.Errorf("expected date 2025-04-21, got %q", o.Date)
	}
}

func TestGetFlights_NoOffers(t *testing.T) {
	handler := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.SerAPIResponse{
			BestFlights: nil,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		data, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		if _, err := w.Write(data); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer handler.Close()

	client := New("key", handler.URL, handler.Client())
	_, err := client.GetFlights(context.Background(), models.FlightSearch{})
	if err == nil || err.Error() != "no flight offers found" {
		t.Fatalf("expected no flight offers found error, got %v", err)
	}
}

func TestGetFlights_HTTPError(t *testing.T) {
	handler := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("something bad")); err != nil {
			t.Fatalf("failed to write error response: %v", err)
		}
	}))
	defer handler.Close()

	client := New("key", handler.URL, handler.Client())
	_, err := client.GetFlights(context.Background(), models.FlightSearch{})
	if err == nil || !strings.Contains(err.Error(), "serpapi error: status 500") {
		t.Fatalf("expected HTTP error from SerpAPI, got %v", err)
	}
}
