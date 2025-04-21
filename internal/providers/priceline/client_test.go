package priceline

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fehepe/flight-price-service/pkg/models"
)

func TestGetFlights_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("originAirportCode") != "AAA" ||
			q.Get("destinationAirportCode") != "BBB" ||
			q.Get("departureDate") != "2025-04-22" {
			t.Fatalf("unexpected query parameters: %v", r.URL.RawQuery)
		}

		w.WriteHeader(http.StatusOK)
		resp := models.PriceLineAPIResponse{
			Data: models.PriceLineData{
				Listings: []models.PriceLineListing{{
					TotalPriceWithDecimal: models.TotalPriceWithDecimal{Price: 123.45},
					Slices: []models.PriceLineSlice{{
						DurationInMinutes: "90",
						Segments: []models.PriceLineSegment{{
							DepartInfo: models.DepartInfo{
								Airport: models.PriceLineAirport{Code: "AAA"},
								Time:    models.PriceLineTime{DateTime: "2025-04-22T10:00:00"},
							},
							ArrivalInfo: models.PriceLineArrivalInfo{Airport: models.PriceLineAirport{Code: "BBB"}},
						}},
					}},
					Airlines: []models.PriceLineAirline{{Name: "AL"}},
				}},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer ts.Close()

	client := New("APIKEY", ts.URL, ts.Client())
	search := models.FlightSearch{
		Origin:        "AAA",
		Destination:   "BBB",
		DepartureDate: time.Date(2025, 4, 22, 0, 0, 0, 0, time.UTC),
	}
	offers, err := client.GetFlights(context.Background(), search)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(offers) != 1 {
		t.Fatalf("expected 1 offer, got %d", len(offers))
	}

	o := offers[0]
	if o.Price != 123.45 {
		t.Errorf("expected price 123.45, got %v", o.Price)
	}
	if o.Origin != "AAA" || o.Destination != "BBB" {
		t.Errorf("expected route AAA->BBB, got %s->%s", o.Origin, o.Destination)
	}
	if o.Date != "2025-04-22" {
		t.Errorf("expected date 2025-04-22, got %s", o.Date)
	}
	if o.Duration != "PT1H30M" {
		t.Errorf("expected duration PT1H30M, got %s", o.Duration)
	}
}

func TestGetFlights_NoFlights(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := models.PriceLineAPIResponse{
			Data: models.PriceLineData{Listings: []models.PriceLineListing{}},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer ts.Close()

	client := New("APIKEY", ts.URL, ts.Client())
	search := models.FlightSearch{Origin: "AAA", Destination: "BBB", DepartureDate: time.Now()}
	offers, err := client.GetFlights(context.Background(), search)
	if err != ErrNoFlights {
		t.Fatalf("expected ErrNoFlights, got %v", err)
	}
	if offers != nil {
		t.Fatalf("expected no offers, got %v", offers)
	}
}

func TestGetFlights_APIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "server error")
	}))
	defer ts.Close()

	client := New("APIKEY", ts.URL, ts.Client())
	search := models.FlightSearch{Origin: "AAA", Destination: "BBB", DepartureDate: time.Now()}
	offers, err := client.GetFlights(context.Background(), search)
	if err == nil || !strings.Contains(err.Error(), "API error") {
		t.Fatalf("expected API error, got %v", err)
	}
	if offers != nil {
		t.Fatalf("expected no offers, got %v", offers)
	}
}

func TestToISO8601(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"90", "PT1H30M"},
		{"61", "PT1H1M"},
		{"0", "PT0H0M"},
		{"invalid", "PT0H0M"},
	}
	for _, tc := range cases {
		out := toISO8601(tc.input)
		if out != tc.expected {
			t.Errorf("toISO8601(%q): expected %q, got %q", tc.input, tc.expected, out)
		}
	}
}
