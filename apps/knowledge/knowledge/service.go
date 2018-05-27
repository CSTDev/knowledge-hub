package knowledge

import (
	"encoding/json"
	"net/http"

	"github.com/cstdev/knowledge-hub/apps/knowledge/database"
	"github.com/cstdev/knowledge-hub/apps/knowledge/types"
	log "github.com/sirupsen/logrus"
)

// WebService provides any dependencies needed by the service
type WebService struct {
	DB database.Database
}

// ErrorResponse is returned for non 200 status'
type ErrorResponse struct {
	Message string
}

// HealthCheck returns if the service is up and running.
func (s *WebService) HealthCheck() http.HandlerFunc {

	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"event": "healthCheck",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// NewRecord adds the passed record to the database
func (s *WebService) NewRecord() http.HandlerFunc {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"event": "create",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var rec types.Record
		decoder := json.NewDecoder(r.Body)

		if r.Body == nil {
			w.WriteHeader(400)
			return
		}

		if err := decoder.Decode(&rec); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to parse JSON"})
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Unable to parse JSON")
			return
		}

		if s.DB == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to connect to database"})
			log.Error("No database set")
			return
		}

		err := s.DB.Create(rec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Failed to create new record"})
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Failed to create new record")
		}

		w.WriteHeader(http.StatusCreated)

	}
}
