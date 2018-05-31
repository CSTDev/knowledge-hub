package types

type Record struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Location location `json:"location"`
	Reports  []Report `json:"reports"`
}

type location struct {
	Lat float64 `json:"lat,string"`
	Lng float64 `json:"lng,string"`
}

type Report struct {
	ReportID      int    `json:"reportID"`
	ReportDetails string `json:"reportDetails"`
	URL           string `json:"url"`
}
