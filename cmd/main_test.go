package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/taofit/coding-challenge-backend/internal"

	"github.com/gorilla/mux"
)

var router *mux.Router

func TestMain(m *testing.M) {
	router = mux.NewRouter()
	initializeRoutes(router)
	code := m.Run()

	os.Exit(code)
}

func TestCreateOfficer(t *testing.T) {
	clearTable()
	dataLoad := []byte(`{"name": "lars mattil"}`)
	req, _ := http.NewRequest("POST", "/officers", bytes.NewBuffer(dataLoad))

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["success"] != "" {
		t.Log(m["success"])
	} else {
		t.Errorf("Expected 'User not created'. Got '%s'", m["error"])
	}

}

func TestGetOfficer(t *testing.T) {
	clearTable()
	addOfficers(1)

	req, _ := http.NewRequest("GET", "/officers/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetNonExistentOfficer(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/officers/74", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["success"] != "" {
		t.Log(m["success"])
	} else {
		t.Errorf("Expected 'User not created'. Got '%s'", m["error"])
	}
}

func TestEmptyOfficers(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/officers", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestUpdateOfficer(t *testing.T) {
	clearTable()
	addOfficers(1)
	req, _ := http.NewRequest("GET", "/officers/1", nil)
	response := executeRequest(req)

	var originalOfficer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalOfficer)
	dataload := []byte(`{"name":"updated officer name"}`)
	req, _ = http.NewRequest("PUT", "/officers/1", bytes.NewBuffer(dataload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	fmt.Printf("%T\n", originalOfficer["id"])
	if m["id"] != fmt.Sprintf("%v", originalOfficer["id"]) {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalOfficer["id"], m["id"])
	}
	if m["name"] == originalOfficer["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalOfficer["name"], m["name"], m["name"])
	}
}

func addOfficers(count int) {
	if count < 1 {
		count = 1
	}
	db := internal.DbConn()
	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO officers(name) VAlUES ('%s')", ("User " + strconv.Itoa(i+1)))
		db.Exec(statement)

	}
	defer db.Close()
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	} else {
		t.Log("NO Error")
	}
}

func clearTable() {
	db := internal.DbConn()

	db.Exec("DELETE FROM officers")
	db.Exec("ALTER TABLE officers AUTO_INCREMENT = 1")

	defer db.Close()
}
