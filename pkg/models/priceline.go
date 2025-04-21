package models

// APIResponse models the top-level JSON structure.
type PriceLineAPIResponse struct {
	Data PriceLineData `json:"data"`
}

type PriceLineData struct {
	Listings []PriceLineListing `json:"listings"`
}

type PriceLineListing struct {
	TotalPriceWithDecimal TotalPriceWithDecimal `json:"totalPriceWithDecimal"`
	Slices                []PriceLineSlice      `json:"slices"`
	Airlines              []PriceLineAirline    `json:"airlines"`
}

type TotalPriceWithDecimal struct {
	Price float64 `json:"price"`
}

type PriceLineSlice struct {
	DurationInMinutes string             `json:"durationInMinutes"`
	Segments          []PriceLineSegment `json:"segments"`
}

type PriceLineSegment struct {
	DepartInfo  DepartInfo           `json:"departInfo"`
	ArrivalInfo PriceLineArrivalInfo `json:"arrivalInfo"`
}

type DepartInfo struct {
	Airport PriceLineAirport `json:"airport"`
	Time    PriceLineTime    `json:"time"`
}

type PriceLineArrivalInfo struct {
	Airport PriceLineAirport `json:"airport"`
}

type PriceLineAirport struct {
	Code string `json:"code"`
}

type PriceLineTime struct {
	DateTime string `json:"dateTime"`
}

type PriceLineAirline struct {
	Name string `json:"name"`
}
