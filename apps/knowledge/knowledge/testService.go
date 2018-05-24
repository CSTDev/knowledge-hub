package knowledge

import (
	"net/http"
)

// TestService example of how to create one for testing
type TestService struct {
}

// HealthCheck for testing
func (s *TestService) HealthCheck() http.HandlerFunc {
	return nil
}
