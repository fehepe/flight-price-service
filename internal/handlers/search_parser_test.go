package handlers

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestExtractFlightSearch(t *testing.T) {
	today := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	tests := []struct {
		name      string
		params    map[string]string
		wantError bool
	}{
		{
			name: "valid request",
			params: map[string]string{
				"origin":      "JFK",
				"destination": "LAX",
				"date":        today,
				"adults":      "2",
				"non_stop":    "true",
			},
			wantError: false,
		},
		{
			name: "missing required fields",
			params: map[string]string{
				"origin": "JFK",
			},
			wantError: true,
		},
		{
			name: "invalid iata",
			params: map[string]string{
				"origin":      "JKF1",
				"destination": "LAX",
				"date":        today,
			},
			wantError: true,
		},
		{
			name: "invalid date",
			params: map[string]string{
				"origin":      "JFK",
				"destination": "LAX",
				"date":        "01-01-2025",
			},
			wantError: true,
		},
		{
			name: "past date",
			params: map[string]string{
				"origin":      "JFK",
				"destination": "LAX",
				"date":        "2000-01-01",
			},
			wantError: true,
		},
		{
			name: "invalid adults",
			params: map[string]string{
				"origin":      "JFK",
				"destination": "LAX",
				"date":        today,
				"adults":      "abc",
			},
			wantError: true,
		},
		{
			name: "too many adults",
			params: map[string]string{
				"origin":      "JFK",
				"destination": "LAX",
				"date":        today,
				"adults":      "9",
			},
			wantError: true,
		},
		{
			name: "invalid non_stop",
			params: map[string]string{
				"origin":      "JFK",
				"destination": "LAX",
				"date":        today,
				"non_stop":    "maybe",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := url.Values{}
			for k, v := range tt.params {
				u.Set(k, v)
			}
			r := &http.Request{URL: &url.URL{RawQuery: u.Encode()}}
			_, err := extractFlightSearch(r)
			if (err != nil) != tt.wantError {
				t.Errorf("extractFlightSearch() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
