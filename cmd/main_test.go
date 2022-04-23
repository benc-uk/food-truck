package main

//
// Integration tests for Trucks API, these make actual HTTP calls, and require the database
//

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/benc-uk/food-truck/pkg/data"
	"github.com/benc-uk/food-truck/pkg/trucks"

	"github.com/gorilla/mux"
)

//
// Setup test server
//
var testServer *httptest.Server

func TestMain(m *testing.M) {
	router := mux.NewRouter()
	truckAPI := &trucks.API{
		Service: trucks.NewService(data.NewDatabase("../data/food-trucks.db")),
	}
	// Add app routes, plus truck service, using a in-memory DB
	truckAPI.AddRoutes(router)

	testServer = httptest.NewServer(router)
	log.Println("Test server listening at:", testServer.URL)
	defer testServer.Close()

	// Now run the rest of the tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

//
// Test cases
//
func TestGetTrucksSanFran(t *testing.T) {
	expected := http.StatusOK
	res, err := http.Get(testServer.URL + "/trucks/37.7758/-122.4205")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != expected {
		t.Errorf("Status code was: %d, expected: %d", res.StatusCode, expected)
	}

	bodyChecker(res, "Snacks", t)
}

func TestGetTrucksLondon(t *testing.T) {
	expected := http.StatusOK
	res, err := http.Get(testServer.URL + "/trucks/51.403278/0.056169")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != expected {
		t.Errorf("Status code was: %d, expected: %d", res.StatusCode, expected)
	}

	bodyChecker(res, "Burgers", t)
}

func TestGetTrucksRadius(t *testing.T) {
	expected := http.StatusOK
	res, err := http.Get(testServer.URL + "/trucks/37.7758/-122.4205?radius=400")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != expected {
		t.Errorf("Status code was: %d, expected: %d", res.StatusCode, expected)
	}

	bodyChecker(res, "Chicken", t)
}

func TestGetTrucksRadiusEmpty(t *testing.T) {
	expected := http.StatusOK
	res, err := http.Get(testServer.URL + "/trucks/0.00/1.23?radius=1000")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != expected {
		t.Errorf("Status code was: %d, expected: %d", res.StatusCode, expected)
	}

	bodyChecker(res, "[]", t)
}

func TestGetTrucksBadParams(t *testing.T) {
	expected := http.StatusBadRequest
	res, err := http.Get(testServer.URL + "/trucks/0.00/foo")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != expected {
		t.Errorf("Status code was: %d, expected: %d", res.StatusCode, expected)
	}
}

func bodyChecker(res *http.Response, expected string, t *testing.T) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	bodyStr := string(body)

	if !strings.Contains(bodyStr, expected) {
		t.Errorf("Payload not as expected, wanted: %s, got: %s", expected, bodyStr)
	}
}
