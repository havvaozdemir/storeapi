package specs

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"storeapi/config"
	"storeapi/data"
	types "storeapi/errors"
	"storeapi/server"
	"storeapi/service"
	"strings"
	"testing"
)

func createTestRequestRecord() data.RequestRecord {
	return data.RequestRecord{
		StartDate: "2016-05-26",
		EndDate:   "2020-02-02",
		MinCount:  2700,
		MaxCount:  3000,
	}
}

func createTestResponseData() data.RequestData {
	return data.RequestData{
		Key:   "active-tabs",
		Value: "getir",
	}
}
func TestFetchRecords(t *testing.T) {
	byteRecordReq, _ := json.Marshal(createTestRequestRecord())
	reader := bytes.NewReader(byteRecordReq)
	req, err := http.NewRequest("POST", "/records", reader)
	if err != nil {
		t.Fatal(err)
	}
	db := server.SetupDB(config.MongoDBURI)
	store := make(map[string]string)
	s := service.New(store, db)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.FetchRecords)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("fetchrecords returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response data.ResponseRecord
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("got invalid response, expected record response, got: %v", rr.Body.String())
	}
	if response.Code != 0 {
		t.Errorf("expected code, got %v", response.Code)
	}
	if response.Msg != "Success" {
		t.Errorf("expected code, got %v", response.Msg)
	}
	if len(response.Records) == 0 {
		t.Error("expected count of records bigger than 0 , got 0")
	}
}

func TestFetchRecordsNotExpectedBody(t *testing.T) {
	reqJSON := `{"sss":"sss"}`

	req, err := http.NewRequest("POST", "/records", strings.NewReader(reqJSON))
	if err != nil {
		t.Fatal(err)
	}
	db := server.SetupDB(config.MongoDBURI)
	store := make(map[string]string)
	s := service.New(store, db)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.FetchRecords)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("fetchrecords returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetStore(t *testing.T) {

	req, err := http.NewRequest("GET", "/in-memory", nil)
	q := req.URL.Query()
	q.Add("key", "active-tabs")
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Fatal(err)
	}
	db := server.SetupDB(config.MongoDBURI)
	store := make(map[string]string)
	store["active-tabs"] = "getir"
	s := service.New(store, db)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetStore)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getStore returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response data.RequestData
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("got invalid response, expected responseData, got: %v", rr.Body.String())
	}
	if response.Value != "getir" {
		t.Errorf("expected key's value, got %v", response.Value)
	}
}

func TestGetStoreNotExistingKey(t *testing.T) {

	req, err := http.NewRequest("GET", "/in-memory", nil)
	q := req.URL.Query()
	q.Add("key", "inActive-tabs")
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Fatal(err)
	}
	db := server.SetupDB(config.MongoDBURI)
	store := make(map[string]string)
	store["active-tabs"] = "getir"
	s := service.New(store, db)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetStore)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("getStore returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGetStoreKeyRequired(t *testing.T) {

	req, err := http.NewRequest("GET", "/in-memory", nil)
	if err != nil {
		t.Fatal(err)
	}
	db := server.SetupDB(config.MongoDBURI)
	store := make(map[string]string)
	store["active-tabs"] = "getir"
	s := service.New(store, db)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetStore)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("getStore returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	expected := `"Key is required!"`
	body := strings.TrimSpace(rr.Body.String())
	if body != expected {
		t.Errorf("getStore returned wrong error message: got %v want %v",
			rr.Body.String(), types.KeyRequired)
	}

}

func TestSetStore(t *testing.T) {
	byteSetStore, _ := json.Marshal(createTestResponseData())
	reader := bytes.NewReader(byteSetStore)
	req, err := http.NewRequest("POST", "/in-memory", reader)
	if err != nil {
		t.Fatal(err)
	}
	db := server.SetupDB(config.MongoDBURI)
	store := make(map[string]string)
	s := service.New(store, db)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.SetStore)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("SetStore returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	var response data.RequestData
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("got invalid response, expected responseData, got: %v", rr.Body.String())
	}
	if response.Value != "getir" {
		t.Errorf("expected key's value getir, got %v", response.Value)
	}
}

func TestSetStoreBodyError(t *testing.T) {
	reqJSON := `{"sss":"sss"}`
	req, err := http.NewRequest("POST", "/in-memory", strings.NewReader(reqJSON))
	if err != nil {
		t.Fatal(err)
	}
	db := server.SetupDB(config.MongoDBURI)
	store := make(map[string]string)
	s := service.New(store, db)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.SetStore)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("SetStore returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected := `"Key is required!"`
	body := strings.TrimSpace(rr.Body.String())
	if body != expected {
		t.Errorf("getStore returned wrong error message: got %v want %v",
			rr.Body.String(), types.KeyRequired)
	}
}
