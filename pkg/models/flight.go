package models

import "time"

// FlightSearch contains search parameters for retrieving flight offers.
type FlightSearch struct {
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	DepartureDate time.Time `json:"departure_date"`
}
