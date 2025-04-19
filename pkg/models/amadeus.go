package models

// AmadeusFlightResponse maps the response from Amadeus flight-offers API
type AmadeusFlightResponse struct {
	Data []AmadeusFlightOffer `json:"data"`
}

type AmadeusFlightOffer struct {
	Itineraries []AmadeusItinerary `json:"itineraries"`
	Price       AmadeusPrice       `json:"price"`
}

type AmadeusPrice struct {
	Total string `json:"total"`
}

type AmadeusItinerary struct {
	Duration string           `json:"duration"`
	Segments []AmadeusSegment `json:"segments"`
}

type AmadeusSegment struct {
	Departure AmadeusLocation `json:"departure"`
	Arrival   AmadeusLocation `json:"arrival"`
	Duration  string          `json:"duration"`
}

type AmadeusLocation struct {
	IataCode string `json:"iataCode"`
	At       string `json:"at"`
}
