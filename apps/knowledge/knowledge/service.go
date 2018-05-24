package knowledge

import (
	"net/http"
)

// Service interface defining all methods expected of a Service
type Service interface {
	HealthCheck() http.HandlerFunc
}

// WebService provides any dependencies needed by the service
type WebService struct {
	DB *string
}

// HealthCheck returns if the service is up and running.
func (s *WebService) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
