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
}

// SearchQuery marshalls the query params into a search term
type SearchQuery struct {
	Query string
}
