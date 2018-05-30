package knowledge

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cstdev/knowledge-hub/apps/knowledge/database"
	"github.com/cstdev/knowledge-hub/apps/knowledge/types"
	"github.com/dyninc/qstring"
	"github.com/gorilla/mux"
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
// Path: /record
// Method: POST
// Example: /record
//		Body: {
//					"title": "A Location",
//					"location": {
//					"lng": "-1.619060757481970",
//					"lat": "53.862309546682600"
//					}
//				}
func (s *WebService) NewRecord() http.HandlerFunc {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"event": "create",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var rec types.Record
		decoder := json.NewDecoder(r.Body)

		if r.Body == nil {
			log.WithFields(log.Fields{
				"status": 400,
			}).Warn("No body provided")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "No body provided"})
			return
		}

		if err := decoder.Decode(&rec); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to parse JSON"})
			log.WithFields(log.Fields{
				"error":  err.Error(),
				"status": 400,
			}).Error("Unable to parse JSON")
			return
		}

		if s.DB == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to connect to database"})
			log.WithFields(log.Fields{
				"status": 500,
			}).Error("No database set")
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

		log.Info("Successfully created new record")
		w.WriteHeader(http.StatusCreated)

	}
}

// Search queries the database for the term passed to the path
// as a URL query parameter.
// Path: /record
// Method: GET
// Parameters: query
// Example: /record?query=Leeds
func (s *WebService) Search() http.HandlerFunc {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"event": "search",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Debug("Called Search")
		query := &types.SearchQuery{}
		log.Debug(r.URL.Query())

		if len(r.URL.Query()) < 1 {
			log.WithFields(log.Fields{
				"status": 400,
			}).Error("No query parameters provided")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "No search parameters provided"})
			return
		}

		err := qstring.Unmarshal(r.URL.Query(), query)

		if err != nil {
			log.WithFields(log.Fields{
				"status": 400,
				"error":  err.Error(),
			}).Error("Unable to unmarshal query params")

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to parse search parameters"})
			return
		}

		if len(query.Query) > 100 {
			log.WithFields(log.Fields{
				"status": 400,
			}).Error("Query string too long")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Query string must be less than 100 characters"})
			return
		}

		records, err := s.DB.Search(*query)
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"status": 500,
				"query":  query,
			}).Error("Unable to search database")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to search database"})
			return
		}

		if len(records) == 0 {
			log.WithFields(log.Fields{
				"status":   404,
				"response": records,
			}).Info("No results found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}

		log.WithFields(log.Fields{
			"status":   200,
			"response": records,
		}).Info("Results returned")
		json.NewEncoder(w).Encode(records)
	}

}

// Update takes a record and writes any changes to the database
// Path: /record
// Method: PUT
// Example: /record/12345
//		Body: {
//					"location": {
//					"lng": "-5.619060757481970",
//					"lat": "52.862309546682600"
//					}
//				}
func (s *WebService) Update() http.HandlerFunc {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"event": "update",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		log.WithFields(log.Fields{
			"vars": vars,
		}).Debug("Passed Vars")

		strID := vars["id"]

		log.WithFields(log.Fields{
			"id": strID,
		}).Debug("Passed Id")

		if strID == "" {
			log.WithFields(log.Fields{
				"status": 400,
			}).Warn("No ID provided")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "No ID provided"})
			return
		}

		id, err := strconv.Atoi(strID)

		if err != nil {
			log.WithFields(log.Fields{
				"status": 400,
				"id":     strID,
			}).Warn("Invalid ID provided")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Invalid ID provided"})
			return
		}

		if r.Body == nil {
			log.WithFields(log.Fields{
				"status": 400,
			}).Warn("No body provided")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "No body provided"})
			return
		}

		var rec types.Record
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&rec); err != nil {
			log.WithFields(log.Fields{
				"error":  err.Error(),
				"status": 400,
			}).Error("Unable to parse JSON")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to parse JSON"})
			return
		}

		if s.DB == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to connect to database"})
			log.WithFields(log.Fields{
				"status": 500,
			}).Error("No database set")
			return
		}

		err = s.DB.Update(id, rec)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Failed to update record"})
			log.WithFields(log.Fields{
				"status": 500,
				"error":  err.Error(),
			}).Error("Failed to update record")
			return
		}
	}
}
