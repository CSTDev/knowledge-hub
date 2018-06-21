package knowledge

import (
	"encoding/json"
	"errors"
	"net/http"

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

type Response struct {
	ID string
}

// HealthCheck returns if the service is up and running.
func (s *WebService) HealthCheck() http.HandlerFunc {

	logger := log.WithFields(log.Fields{
		"event": "healthCheck",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		logger.Info("Health Check OK")
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
	logger := log.WithFields(log.Fields{
		"event": "create",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var rec types.Record
		decoder := json.NewDecoder(r.Body)

		if r.Body == nil {
			logger.WithFields(log.Fields{
				"status": 400,
			}).Warn("No body provided")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "No body provided"})
			return
		}

		if err := decoder.Decode(&rec); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to parse JSON"})
			logger.WithFields(log.Fields{
				"error":  err.Error(),
				"status": 400,
			}).Error("Unable to parse JSON")
			return
		}

		if s.DB == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to connect to database"})
			logger.WithFields(log.Fields{
				"status": 500,
			}).Error("No database set")
			return
		}

		id, err := s.DB.Create(rec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Failed to create new record"})
			logger.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Failed to create new record")
			return
		}

		logger.WithFields(log.Fields{
			"id":     id,
			"status": 201,
		}).Info("Successfully created new record")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&Response{ID: id})
	}
}

// Search queries the database for the term passed to the path
// as a URL query parameter.
// Path: /record
// Method: GET
// Parameters:
//		query - string to search for
//		minLat - minimum latitude of the location
//		maxLat - maximum latitude of the location
//		minLng - mininum longitude of the location
//		maxLng - maximum longitude of the location
// Example: /record?query=Leeds&minLat=-1.23423423&maxLat=0.12321321&minLng=54.4564523&maxLng=55.2342809
func (s *WebService) Search() http.HandlerFunc {
	logger := log.WithFields(log.Fields{
		"event": "search",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		logger.Debug("Called Search")
		query := &types.SearchQuery{}
		logger.Debug(r.URL.Query())

		if len(r.URL.Query()) < 1 {
			logger.WithFields(log.Fields{
				"status": 400,
			}).Error("No query parameters provided")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "No search parameters provided"})
			return
		}

		err := qstring.Unmarshal(r.URL.Query(), query)

		logger.WithFields(log.Fields{
			"value": query,
		}).Debug("Unmarshalled query params")

		if err != nil {
			logger.WithFields(log.Fields{
				"status": 400,
				"error":  err.Error(),
			}).Error("Unable to unmarshal query params")

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to parse search parameters"})
			return
		}

		if !boundsPresent(*query) {
			logger.WithFields(log.Fields{
				"status": 400,
			}).Error("No location bounds provided")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Not all location bounds provided, expected mininum and maximum latitude and longitude."})
			return
		}

		if len(query.Query) > 100 {
			logger.WithFields(log.Fields{
				"status": 400,
			}).Error("Query string too long")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Query string must be less than 100 characters"})
			return
		}

		records, err := s.DB.Search(*query)
		if err != nil {
			logger.WithFields(log.Fields{
				"error":  err,
				"status": 500,
				"query":  query,
			}).Error("Unable to search database")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to search database"})
			return
		}

		if len(records) == 0 {
			logger.WithFields(log.Fields{
				"status":   404,
				"response": records,
			}).Info("No results found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}

		logger.WithFields(log.Fields{
			"status":         200,
			"numberReturned": len(records),
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
	logger := log.WithFields(log.Fields{
		"event": "update",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, err := getRecordID(r)
		if err != nil {
			logger.WithFields(log.Fields{
				"status": 400,
				"error":  err.Error(),
			}).Warn("Issue with ID")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: err.Error()})
			return
		}

		if r.Body == nil {
			logger.WithFields(log.Fields{
				"status": 400,
			}).Warn("No body provided")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "No body provided"})
			return
		}

		var rec types.Record
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&rec); err != nil {
			logger.WithFields(log.Fields{
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
			logger.WithFields(log.Fields{
				"status": 500,
			}).Error("No database set")
			return
		}

		rec.ID = id

		err = s.DB.Update(id, rec)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Failed to update record"})
			logger.WithFields(log.Fields{
				"status": 500,
				"error":  err.Error(),
			}).Error("Failed to update record")
			return
		}

		logger.WithFields(log.Fields{
			"status": 200,
			"id":     id,
		}).Info("Updated record")
		w.WriteHeader(http.StatusOK)
	}
}

// Delete takes a record ID and marks it as deleted
// Path: /record
// Method: DELETE
// Example: /record/12345
func (s *WebService) Delete() http.HandlerFunc {
	logger := log.WithFields(log.Fields{
		"event": "delete",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, err := getRecordID(r)
		if err != nil {
			logger.WithFields(log.Fields{
				"status": 400,
				"error":  err.Error(),
			}).Warn("Issue with ID")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: err.Error()})
			return
		}

		err = s.DB.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Failed to delete record"})
			logger.WithFields(log.Fields{
				"status": 500,
				"error":  err.Error(),
			}).Error("Failed to delete record")
			return
		}

		logger.WithFields(log.Fields{
			"status": 200,
			"id":     id,
		}).Info("Deleted record")
		w.WriteHeader(http.StatusOK)
	}

}

// GetFields retrieves the fields that the user can enter from the database
// Path: /field
// Method: GET
// Example: /field
func (s *WebService) GetFields() http.HandlerFunc {

	logger := log.WithFields(log.Fields{
		"event": "GetFields",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fields, err := s.DB.Fields()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Failed to get fields"})
			logger.WithFields(log.Fields{
				"status": 500,
				"error":  err.Error(),
			}).Error("Failed to get fields")
			return
		}

		if len(fields) == 0 {
			logger.WithFields(log.Fields{
				"status": 404,
			}).Info("No fields found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("[]"))
			return
		}

		logger.WithFields(log.Fields{
			"fieldCount": len(fields),
		}).Info("Returning fields")
		json.NewEncoder(w).Encode(fields)

	}

}

// UpdateFields stores the fields set in the database
// Path: /field
// Method: PUT
// Example: /field
//		Body: [
//				{
//					"id": "123456",
//					"value": "Field Name",
//					"order": 0
//				}
//			]
func (s *WebService) UpdateFields() http.HandlerFunc {
	logger := log.WithFields(log.Fields{
		"event": "UpdateFields",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Body == nil {
			logger.WithFields(log.Fields{
				"status": 400,
			}).Warn("No body provided")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "No body provided"})
			return
		}

		var fields []types.Field
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&fields); err != nil {
			logger.WithFields(log.Fields{
				"error":  err.Error(),
				"status": 400,
			}).Error("Unable to parse JSON")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to parse JSON"})
			return
		}

		err := s.DB.UpdateFields(fields)
		if err != nil {
			logger.WithFields(log.Fields{
				"error":  err.Error(),
				"status": 500,
			}).Error("Unable to store fields")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: "Unable to store fields"})
			return
		}

		logger.WithFields(log.Fields{
			"status": 200,
		}).Info("Updated Fields")

	}
}

// DeleteField takes a field id and marks it as deleted
// Path: /field
// Method: DELETE
// Example: /field/12345
func (s *WebService) DeleteField() http.HandlerFunc {
	logger := log.WithFields(log.Fields{
		"event": "deleteField",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, err := getRecordID(r)
		if err != nil {
			logger.WithFields(log.Fields{
				"status": 400,
				"error":  err.Error(),
			}).Warn("Issue with ID")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: err.Error()})
			return
		}

		err = s.DB.DeleteField(id)
		if err != nil {

			_, ok := err.(*types.FieldNotFoundError)
			logger.WithFields(log.Fields{
				"result": ok,
			}).Debug("Check if error from database is FieldNotFound")

			if ok {
				logger.WithFields(log.Fields{
					"status": 404,
					"error":  err.Error(),
					"id":     id,
				}).Warn("Couldn't find field to delete")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&ErrorResponse{Message: err.Error()})
				return
			}

			logger.WithFields(log.Fields{
				"status": 500,
				"error":  err.Error(),
				"id":     id,
			}).Error("Failed to delete field")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Message: err.Error()})
			return
		}

		logger.WithFields(log.Fields{
			"status": 200,
			"id":     id,
		}).Info("Deleted field")
		w.WriteHeader(http.StatusOK)
	}

}

func getRecordID(r *http.Request) (string, error) {
	vars := mux.Vars(r)

	log.WithFields(log.Fields{
		"vars": vars,
	}).Debug("Passed Vars")

	strID := vars["id"]

	log.WithFields(log.Fields{
		"id": strID,
	}).Debug("Passed Id")

	if strID == "" {
		return "", errors.New("No ID provided")
	}

	// id, err := strconv.Atoi(strID)

	// if err != nil {
	// 	return 0, errors.New("Invalid ID provided")
	// }

	return strID, nil
}

func boundsPresent(query types.SearchQuery) bool {
	if query.MinLat >= query.MaxLat || query.MinLng >= query.MaxLng {
		return false
	}
	return query.MinLat != 0 && query.MaxLat != 0 && query.MinLng != 0 && query.MaxLng != 0
}
