package models

// APIResponse models the top-level JSON structure.
type PriceLineAPIResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Listings []Listing `json:"listings"`
}

type Listing struct {
	TotalPriceWithDecimal TotalPriceWithDecimal `json:"totalPriceWithDecimal"`
	Slices                []Slice               `json:"slices"`
	Airlines              []Airline             `json:"airlines"`
}

type TotalPriceWithDecimal struct {
	Price float64 `json:"price"`
}

type Slice struct {
	DurationInMinutes string    `json:"durationInMinutes"`
	Segments          []Segment `json:"segments"`
}

type Segment struct {
	DepartInfo  DepartInfo  `json:"departInfo"`
	ArrivalInfo ArrivalInfo `json:"arrivalInfo"`
}

type DepartInfo struct {
	Airport Airport `json:"airport"`
	Time    Time    `json:"time"`
}

type ArrivalInfo struct {
	Airport Airport `json:"airport"`
}

type Airport struct {
	Code string `json:"code"`
}

type Time struct {
	DateTime string `json:"dateTime"`
}

type Airline struct {
	Name string `json:"name"`
}
