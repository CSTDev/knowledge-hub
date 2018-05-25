package knowledge

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cstdev/knowledge-hub/apps/knowledge/database"
	"github.com/cstdev/knowledge-hub/apps/knowledge/types"
)

// WebService provides any dependencies needed by the service
type WebService struct {
	DB database.Database
}

// HealthCheck returns if the service is up and running.
func (s *WebService) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// NewRecord adds the passed record to the database
func (s *WebService) NewRecord() http.HandlerFunc {
	fmt.Println("In")
	return func(w http.ResponseWriter, r *http.Request) {
		var rec types.Record
		decoder := json.NewDecoder(r.Body)

		if r.Body == nil {
			w.WriteHeader(400)
			return
		}

		if err := decoder.Decode(&rec); err != nil {
			w.WriteHeader(422)
			return
		}
		err := s.DB.Create(rec)

		if err != nil {
			fmt.Println(rec)
		}

	}
}
