package models

type FlightOffer struct {
	Provider    string  `json:"provider"`
	Price       float64 `json:"price"`
	Duration    string  `json:"duration"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Date        string  `json:"date"`
}

type SearchResponse struct {
	Cheapest  FlightOffer              `json:"cheapest"`
	Fastest   FlightOffer              `json:"fastest"`
	Providers map[string][]FlightOffer `json:"providers"`
}
