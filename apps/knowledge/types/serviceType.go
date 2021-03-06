package types

import "net/http"

// Service interface defining all methods expected of a Service
type Service interface {
	HealthCheck() http.HandlerFunc
	NewRecord() http.HandlerFunc
	Search() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	GetFields() http.HandlerFunc
	UpdateFields() http.HandlerFunc
	DeleteField() http.HandlerFunc
}

// SearchQuery marshalls the query params into a search term
type SearchQuery struct {
	Query  string
	MinLat float64 `qstring:"minLat"`
	MaxLat float64 `qstring:"maxLat"`
	MinLng float64 `qstring:"minLng"`
	MaxLng float64 `qstring:"maxLng"`
}
