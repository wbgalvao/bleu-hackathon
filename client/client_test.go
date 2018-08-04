package client

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

var (
	baseURL = "https://bleutrade.com/api/v2/"
	apiKey  = "dcf7c6651d91ab69c31df6355f9ea8e28f4ab4ed" // ReadOnly access
	c       Client
)

func setup() {
	// setup test client
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("error parsing base URL: %v\n", err)
	}
	c.BaseURL = u
	c.APIKey = apiKey
	c.httpClient = new(http.Client)
}

func TestBuildRequest(t *testing.T) {
	_, err := c.BuildRequest("GET", "/public/getcurrencies", nil)
	if err != nil {
		t.Errorf("error building HTTP request: %v\n", err)
	}
}

func TestGetBalances(t *testing.T) {
	_, err := c.GetBalances()
	if err != nil {
		t.Errorf("error calling GetBalances endpoint: %v\n", err)
	}
}

// func TestGetBalances_checkResults(t *testing.T) {
// 	results, _ := c.GetBalances()
// 	if len(results) == 0 {
// 		t.Errorf("could not retrieve any balance from GetBalances endpoint\n")
// 	}
// }

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}
