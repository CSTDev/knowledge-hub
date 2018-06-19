package knowledge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cstdev/knowledge-hub/apps/knowledge/types"
	"github.com/gorilla/mux"
)

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

type mockDB struct {
	SearchQuery      types.SearchQuery
	CreateFunc       func(r types.Record) (string, error)
	SearchFunc       func(query types.SearchQuery) ([]types.Record, error)
	UpdateFunc       func(id string, r types.Record) error
	DeleteFunc       func(id string) error
	GetFieldsFunc    func() ([]types.Field, error)
	UpdateFieldsFunc func(f []types.Field) error
}

func (db *mockDB) Create(r types.Record) (string, error) {
	return db.CreateFunc(r)
}

func (db *mockDB) Search(query types.SearchQuery) ([]types.Record, error) {
	return db.SearchFunc(query)
}

func (db *mockDB) Update(id string, r types.Record) error {
	return db.UpdateFunc(id, r)
}

func (db *mockDB) Delete(id string) error {
	return db.DeleteFunc(id)
}

func (db *mockDB) Fields() ([]types.Field, error) {
	return db.GetFieldsFunc()
}

func (db *mockDB) UpdateFields(f []types.Field) error {
	return db.UpdateFieldsFunc(f)
}

var called bool

var jsonReq = []byte(`{
	"title": "Holy Trinity Church",
	"location": {
	  "lng": -1.619060757481970,
	  "lat": 53.862309546682600
	},
	"reports": [
	  {
		"reportID": 0,
		"reportDetails": "that lightsaber times, by but star consists ",
		"url": "https://example.edu/"
	  },
	  {
		"reportID": 1,
		"reportDetails": "can the or brightly burn, fictional length).[1] ",
		"url": "http://example.com/#bat"
	  }
	 ]
}`)

func TestNewRecordCallsDatabaseWithRecord(t *testing.T) {
	called = false
	db := mockDB{
		CreateFunc: func(r types.Record) (string, error) {
			called = true
			return "ABC123", nil
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("POST", "/record", bytes.NewBuffer(jsonReq))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.NewRecord())
	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("Expected database create method to be called")
	}
}

func TestReturnsErrorOnEmptyBody(t *testing.T) {
	db := mockDB{}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("POST", "/record", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.NewRecord())
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Error("Expected 400 status to be returned")
	}
}

func TestReturnsErrorOnInvalidJson(t *testing.T) {
	db := mockDB{}
	service := &WebService{DB: &db}

	badJSON := []byte(`{"abc":123, def:123}`)

	req, err := http.NewRequest("POST", "/record", bytes.NewBuffer(badJSON))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.NewRecord())
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Error("Expected 400 status to be returned")
	}
}

func TestReturnsServerErrorWhenNoDB(t *testing.T) {
	service := &WebService{}
	req, err := http.NewRequest("POST", "/record", bytes.NewBuffer(jsonReq))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.NewRecord())
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected internal server error status to be returned got %d", rr.Code)
	}
}

func TestReturnsCreatedStatusOnSuccess(t *testing.T) {
	db := mockDB{
		CreateFunc: func(r types.Record) (string, error) {
			return "ABC123", nil
		},
	}
	service := &WebService{DB: &db}
	req, err := http.NewRequest("POST", "/record", bytes.NewBuffer(jsonReq))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.NewRecord())
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected Created (201) status to be returned got %d", rr.Code)
	}
}

func TestSearchWithNoQueryParamsErrors(t *testing.T) {
	db := mockDB{}
	service := &WebService{DB: &db}
	req, err := http.NewRequest("GET", "/record", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Search())
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected Bad Request (400) status to be returned got %d", rr.Code)
	}
}

func TestSearchQueryTooLong(t *testing.T) {
	db := mockDB{}
	service := &WebService{DB: &db}
	req, err := http.NewRequest("GET", "/record?query=A%20Really%20long%20query%20A%20Really%20long%20query%20A%20Really%20long%20query%20A%20Really%20long%20query%20A%20Really%20long%20query%20A%20Really%20long%20query%20", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Search())
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected Bad Request (400) status to be returned got %d", rr.Code)
	}
}

func TestSearchQueryIsPassedToDB(t *testing.T) {
	var passedQuery string
	db := mockDB{
		SearchFunc: func(s types.SearchQuery) ([]types.Record, error) {
			passedQuery = s.Query
			return []types.Record{}, nil
		},
	}
	service := &WebService{DB: &db}

	expectedQuery := "Leeds"

	req, err := http.NewRequest("GET", "/record?query=Leeds", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Search())
	handler.ServeHTTP(rr, req)
	if passedQuery != expectedQuery {
		t.Error("Expected search passed to DB but it wasn't")
	}
}

func TestSearchReturnsServerErrorWhenDBSearchFails(t *testing.T) {
	db := mockDB{
		SearchFunc: func(query types.SearchQuery) ([]types.Record, error) {
			return []types.Record{}, errors.New("Unable to search")
		},
	}
	service := &WebService{DB: &db}
	req, err := http.NewRequest("GET", "/record?query=Leeds", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Search())
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected Internal Server Error (500) status to be returned got %d", rr.Code)
	}
}

func TestSuccessfulSearchReturnsResults(t *testing.T) {
	expectedResults := `[{"id":"","title":"Holy Trinity Church","location":{"lat":53.8623095466826,"lng":-1.61906075748197},"reports":[{"reportID":0,"reportDetails":"that lightsaber times, by but star consists ","url":"https://example.edu/"}]}]`
	db := mockDB{
		SearchFunc: func(query types.SearchQuery) ([]types.Record, error) {
			var records []types.Record
			buf := bytes.NewBuffer([]byte(expectedResults))
			err := json.NewDecoder(buf).Decode(&records)
			ok(t, err)
			return records, nil
		},
	}
	service := &WebService{DB: &db}
	req, err := http.NewRequest("GET", "/record?query=Leeds", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Search())
	handler.ServeHTTP(rr, req)
	if strings.TrimSpace(rr.Body.String()) != expectedResults {
		t.Errorf("Expected response to be: \n %s \n but got: \n %s", expectedResults, rr.Body.String())
		t.FailNow()
	}

	if rr.Code != 200 {
		t.Errorf("Expected Success (200) status to be returned got %d", rr.Code)
	}

}

func TestSearchReturnsNotFoundOnNoResults(t *testing.T) {
	db := mockDB{
		SearchFunc: func(query types.SearchQuery) ([]types.Record, error) {
			var records []types.Record
			return records, nil
		},
	}
	service := &WebService{DB: &db}
	req, err := http.NewRequest("GET", "/record?query=Leeds", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Search())
	handler.ServeHTTP(rr, req)
	if strings.TrimSpace(rr.Body.String()) != "{}" {
		t.Errorf("Expected response to be: \n {} \n but got: \n %s", rr.Body.String())
		t.FailNow()
	}

	if rr.Code != 404 {
		t.Errorf("Expected Success (404) status to be returned got %d", rr.Code)
	}
}

func TestUpdateReturnsErrorOnEmptyBody(t *testing.T) {
	db := mockDB{}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("PUT", "/record/12345", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Update())
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Errorf("Expected Bad Request (400) status to be returned got %d", rr.Code)
	}
}

func TestUpdateReturnsErrorOnInvalidJSON(t *testing.T) {
	db := mockDB{}
	service := &WebService{DB: &db}

	badJSON := []byte(`{"abc":123, def:123}`)
	req, err := http.NewRequest("PUT", "/record/12345", bytes.NewBuffer(badJSON))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Update())
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Errorf("Expected Bad Request (400) status to be returned got %d", rr.Code)
	}
}

func TestUpdateReturnsServerErrorWhenNoDB(t *testing.T) {
	service := &WebService{}
	req, err := http.NewRequest("PUT", "/record/12345", bytes.NewBuffer(jsonReq))
	ok(t, err)

	rr := httptest.NewRecorder()
	updateRouter(service).ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected internal server error status to be returned got %d", rr.Code)
	}
}

func TestUpdateReturnsErrorIfNoIdIsProvided(t *testing.T) {

	service := &WebService{}

	req, err := http.NewRequest("PUT", "/record", bytes.NewBuffer(jsonReq))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Update())
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Errorf("Expected Bad Request (400) status to be returned got %d", rr.Code)
	}
}

func updateRouter(service *WebService) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/record/{id}", service.Update())
	return r
}

func TestUpdateCallsDatabaseWithRecordAndId(t *testing.T) {
	called = false
	var passedID string
	db := mockDB{
		UpdateFunc: func(id string, r types.Record) error {
			called = true
			passedID = id
			return nil
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("PUT", "/record/12345", bytes.NewBuffer(jsonReq))
	ok(t, err)

	rr := httptest.NewRecorder()
	updateRouter(service).ServeHTTP(rr, req)

	if !called {
		t.Error("Expected database update method to be called")
		t.FailNow()
	}

	if passedID != "12345" {
		t.Errorf("Expected id: %s \n Got Id: %s \n", "12345", passedID)
	}
}

func TestUpdateReturnsErrorWhenDBUpdateFails(t *testing.T) {

	db := mockDB{
		UpdateFunc: func(id string, r types.Record) error {
			return errors.New("Database failed")
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("PUT", "/record/12345", bytes.NewBuffer(jsonReq))
	ok(t, err)

	rr := httptest.NewRecorder()
	updateRouter(service).ServeHTTP(rr, req)

	if rr.Code != 500 {
		t.Errorf("Expected Internal Server Error (500) status to be returned got %d", rr.Code)
	}
}

func TestUpdateReturnsOkIfRecordIsUpdated(t *testing.T) {
	db := mockDB{
		UpdateFunc: func(id string, r types.Record) error {
			return nil
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("PUT", "/record/12345", bytes.NewBuffer(jsonReq))
	ok(t, err)

	rr := httptest.NewRecorder()
	updateRouter(service).ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected OK (200) status to be returned got %d", rr.Code)
	}
}

func TestDeleteReturnsErrorIfNoIdProvided(t *testing.T) {
	service := &WebService{}

	req, err := http.NewRequest("DELETE", "/record", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.Delete())
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Errorf("Expected Bad Request (400) status to be returned got %d", rr.Code)
	}
}

func deleteRouter(service *WebService) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/record/{id}", service.Delete())
	return r
}

func TestDeleteCallsDatabaseDelete(t *testing.T) {
	called = false
	var passedID string
	db := mockDB{
		DeleteFunc: func(id string) error {
			called = true
			passedID = id
			return nil
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("DELETE", "/record/12345", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	deleteRouter(service).ServeHTTP(rr, req)

	if !called {
		t.Error("Expected database delete method to be called")
		t.FailNow()
	}

	expectedID := "12345"
	if passedID != expectedID {
		t.Errorf("Expected id: %s \n Got Id: %s \n", expectedID, passedID)
	}
}

func TestOnDeleteErrorServerErrorIsReturned(t *testing.T) {
	db := mockDB{
		DeleteFunc: func(id string) error {
			return errors.New("Database failed")
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("DELETE", "/record/12345", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	deleteRouter(service).ServeHTTP(rr, req)

	if rr.Code != 500 {
		t.Errorf("Expected Internal Server Error (500) status to be returned got %d", rr.Code)
	}
}

func TestSuccessfulDeleteOkIsReturned(t *testing.T) {
	db := mockDB{
		DeleteFunc: func(id string) error {
			return nil
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("DELETE", "/record/12345", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	deleteRouter(service).ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected Ok (200) status to be returned got %d", rr.Code)
	}
}

func TestGetFieldsCallsGetFieldsOnTheDatabase(t *testing.T) {
	called = false

	db := mockDB{
		GetFieldsFunc: func() ([]types.Field, error) {
			called = true
			return nil, nil
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("GET", "/field", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.GetFields())
	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("Expected database Get Fields method to be called")
		t.FailNow()
	}
}

func TestGetFieldsReturnsStatus500WhenFieldsCannotBeRetreived(t *testing.T) {
	db := mockDB{
		GetFieldsFunc: func() ([]types.Field, error) {
			return nil, errors.New("Failed to get fields")
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("GET", "/field", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.GetFields())
	handler.ServeHTTP(rr, req)

	if rr.Code != 500 {
		t.Error("Expected status 500 to be returned")
		t.FailNow()
	}
}

func TestGetFieldsReturnsAJSONArrayOfFields(t *testing.T) {
	db := mockDB{
		GetFieldsFunc: func() ([]types.Field, error) {
			fields := []types.Field{
				types.Field{
					ID:    "1",
					Value: "Question 1",
					Order: 1,
				},
				types.Field{
					ID:    "2",
					Value: "Question 2",
					Order: 2,
				},
			}
			return fields, nil
		},
	}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("GET", "/field", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.GetFields())
	handler.ServeHTTP(rr, req)

	expectedJSON := `[{"id":"1","value":"Question 1","order":1},{"id":"2","value":"Question 2","order":2}]`
	responseJSON := strings.TrimSpace(rr.Body.String())
	if expectedJSON != responseJSON {
		t.Errorf("Response JSON didn't match expected. \n Expected: %s \n Got: %s", expectedJSON, responseJSON)
		t.FailNow()
	}
}

func TestWhenNoFieldsAreFoundAnEmptyArrayIsReturned(t *testing.T) {
	db := mockDB{
		GetFieldsFunc: func() ([]types.Field, error) {
			var fields []types.Field
			return fields, nil
		},
	}
	service := &WebService{DB: &db}
	req, err := http.NewRequest("GET", "/field", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.GetFields())
	handler.ServeHTTP(rr, req)
	if strings.TrimSpace(rr.Body.String()) != "[]" {
		t.Errorf("Expected response to be: \n [] \n but got: \n %s", rr.Body.String())
		t.FailNow()
	}

	if rr.Code != 404 {
		t.Errorf("Expected Success (404) status to be returned got %d", rr.Code)
	}
}

func TestUpdateFieldsReturnsErrorOnEmptyBody(t *testing.T) {
	db := mockDB{}
	service := &WebService{DB: &db}

	req, err := http.NewRequest("PUT", "/field", nil)
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.UpdateFields())
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Errorf("Expected Bad Request (400) status to be returned got %d", rr.Code)
	}
}

func TestUpdateFieldReturnsErrorOnInvalidJSONBody(t *testing.T) {
	db := mockDB{}
	service := &WebService{DB: &db}

	badJSON := []byte(`{"abc":123, def:123}`)
	req, err := http.NewRequest("PUT", "/field", bytes.NewBuffer(badJSON))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.UpdateFields())
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Errorf("Expected Bad Request (400) status to be returned got %d", rr.Code)
	}
}

func TestUpdateFieldCallsDBWithFields(t *testing.T) {
	called = false
	db := mockDB{
		UpdateFieldsFunc: func(f []types.Field) error {
			called = true
			return nil
		},
	}
	service := &WebService{DB: &db}

	fieldJSON := []byte(`[{"id":"1","value":"Question 1","order":1},{"id":"2","value":"Question 2","order":2}]`)

	req, err := http.NewRequest("PUT", "/field", bytes.NewBuffer(fieldJSON))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.UpdateFields())
	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("Expected database create method to be called")
	}
}

func TestStatus500ReturnedIfUnableToWriteFieldsToDatabase(t *testing.T) {

	db := mockDB{
		UpdateFieldsFunc: func(f []types.Field) error {
			return errors.New("Failed to write to DB")
		},
	}
	service := &WebService{DB: &db}

	fieldJSON := []byte(`[{"id":"1","value":"Question 1","order":1},{"id":"2","value":"Question 2","order":2}]`)

	req, err := http.NewRequest("PUT", "/field", bytes.NewBuffer(fieldJSON))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.UpdateFields())
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected internal server error status to be returned got %d", rr.Code)
	}

}

func TestSuccessfulUpdateOfFieldsReturns200Status(t *testing.T) {
	db := mockDB{
		UpdateFieldsFunc: func(f []types.Field) error {
			return nil
		},
	}
	service := &WebService{DB: &db}

	fieldJSON := []byte(`[{"id":"1","value":"Question 1","order":1},{"id":"2","value":"Question 2","order":2}]`)

	req, err := http.NewRequest("PUT", "/field", bytes.NewBuffer(fieldJSON))
	ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.UpdateFields())
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected OK status to be returned got %d", rr.Code)
	}
}
