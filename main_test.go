package main

import (
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM endpoints")
	a.DB.Exec("ALTER SEQUENCE endpoints_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS endpoints
(
    id SERIAL,
    url TEXT NOT NULL,
		content TEXT NOT NULL DEFAULT 0.00,
		CONSTRAINT endpoints_pkey PRIMARY KEY (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/endpoints", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentEndpoint(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/endpoint/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestCreateEndpoint(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"url":"testendpoint", "content": "{}"}`)
	req, _ := http.NewRequest("POST", "/api/endpoint", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["url"] != "testendpoint" {
		t.Errorf("Expected endpoint url to be 'testendpoint'. Got '%v'", m["url"])
	}

	if m["content"] != "{}" {
		t.Errorf("Expected endpoint content to be '{}'. Got '%v'", m["content"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetEndpoint(t *testing.T) {
	clearTable()
	addEndpoints(1)

	req, _ := http.NewRequest("GET", "/api/endpoint/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addEndpoints(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO endpoints(url, content) VALUES($1, $2)", "endpoint"+strconv.Itoa(i), strconv.Itoa((i+1.0)*10))
	}
}

func TestUpdateEndpoint(t *testing.T) {

	clearTable()
	addEndpoints(1)

	req, _ := http.NewRequest("GET", "/api/endpoint/1", nil)
	response := executeRequest(req)
	var originalEndpoint map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalEndpoint)

	var jsonStr = []byte(`{"url":"endpoint-updated", "content": "{}"}`)
	req, _ = http.NewRequest("PUT", "/api/endpoint/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalEndpoint["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalEndpoint["id"], m["id"])
	}

	if m["url"] == originalEndpoint["url"] {
		t.Errorf("Expected the url to change from '%v' to '%v'. Got '%v'", originalEndpoint["url"], m["url"], m["url"])
	}

	if m["content"] == originalEndpoint["content"] {
		t.Errorf("Expected the content to change from '%v' to '%v'. Got '%v'", originalEndpoint["content"], m["content"], m["content"])
	}
}

func TestDeleteEndpoint(t *testing.T) {
	clearTable()
	addEndpoints(1)

	req, _ := http.NewRequest("GET", "/api/endpoint/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/endpoint/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/endpoint/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
