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

func TestNewRecordCallsDatabaseWithRecord(t *testing.T) {
	called = false
	db := mockDB{}
	service := &WebService{DB: &db}

	jsonReq := []byte(`{
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

}
