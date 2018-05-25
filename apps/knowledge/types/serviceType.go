package types

import "net/http"

// Service interface defining all methods expected of a Service
type Service interface {
	HealthCheck() http.HandlerFunc
	NewRecord() http.HandlerFunc
}
