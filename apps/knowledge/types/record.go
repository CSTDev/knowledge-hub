package types

type Record struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Location   location               `json:"location"`
	ShortName  string                 `json:"shortName"`
	Facilities []string               `json:"facilities"`
	Details    map[string]interface{} `json:"details"`
}

type location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	Country     string    `json:"country"`
}
