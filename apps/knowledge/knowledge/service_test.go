package knowledge

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cstdev/knowledge-hub/apps/knowledge/types"
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
}

func (db *mockDB) Create(r types.Record) error {
	called = true
	return nil
}

var called bool

var jsonReq = []byte(`{
	"title": "Holy Trinity Church",
	"location": {
	  "lng": "-1.619060757481970",
	  "lat": "53.862309546682600"
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
	db := mockDB{}
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
	db := mockDB{}
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
