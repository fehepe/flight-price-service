package models

// top‐level container
type SerAPIResponse struct {
	BestFlights []FlightOption `json:"best_flights"`
}

type FlightOption struct {
	Flights       []FlightSegment `json:"flights"`
	TotalDuration int             `json:"total_duration"`
	Price         int             `json:"price"`
}

type AirportInfo struct {
	ID   string `json:"id"`
	Time string `json:"time"`
}

type FlightSegment struct {
	DepartureAirport AirportInfo `json:"departure_airport"`
	ArrivalAirport   AirportInfo `json:"arrival_airport"`
	Duration         int         `json:"duration"`
}
