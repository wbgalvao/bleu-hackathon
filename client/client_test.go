package client

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
)

var (
	baseURL   = "https://bleutrade.com/api/v2/"
	apiKey    = "af849b626359ccf9a13f0d8fde899c97" // ReadOnly access
	apiSecret = "ee8f3d52cadcb93b00b96f0be3d7d135a206e4ad"
	c         Client
)

func setup() {
	// setup test client
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("error parsing base URL: %v\n", err)
	}
	c.BaseURL = u
	c.APIKey = apiKey
	c.APISecret = apiSecret
	c.httpClient = new(http.Client)
}

func TestBuildRequest(t *testing.T) {
	_, err := c.BuildRequest("GET", "/public/getcurrencies", nil)
	if err != nil {
		t.Errorf("error building HTTP request: %v\n", err)
	}

}

func TestGetBalances(t *testing.T) {
	balances, err := c.GetBalances()
	t.Error("[LOG]")
	if err != nil {
		t.Errorf("error calling GetBalances endpoint: %v\n", err)
	}

	for _, balance := range balances {
		if n, err := strconv.ParseFloat(balance.Available, 32); n > 0 && err == nil {
			t.Logf("[LOG] Coin: %s | Value: %s\n", balance.Currency, balance.Available)
		}
	}
}

func TestGetBalancesOptionalParameter(t *testing.T) {
	balances, err := c.GetBalances("DOGE")
	t.Error("[LOG]")
	if err != nil {
		t.Errorf("error calling GetBalances endpoint: %v\n", err)
	}

	for _, balance := range balances {
		if n, err := strconv.ParseFloat(balance.Available, 32); n > 0 && err == nil {
			t.Logf("[LOG] Coin: %s | Value: %s\n", balance.Currency, balance.Available)
		}
	}
}

func TestListOrders(t *testing.T) {
	orders, err := c.ListOrder("ALL", "ALL", "ALL")
	t.Error("[LOG]")
	if err != nil {
		t.Errorf("error calling GetBalances endpoint: %v\n", err)
	}

	for _, order := range orders {
		t.Logf("[LOG] Type: %s | Status: %s|Created at: %s| Price: %s\n", order.Type, order.Status, order.Created, order.Price)
	}
}

func TestListOrdersOptionalParameter(t *testing.T) {
	orders, err := c.ListOrder("ALL", "ALL", "ALL", "1000")
	t.Error("[LOG]")
	if err != nil {
		t.Errorf("error calling GetBalances endpoint: %v\n", err)
	}

	for _, order := range orders {
		t.Logf("[LOG] Type: %s | Status: %s|Created at: %s| Price: %s\n", order.Type, order.Status, order.Created, order.Price)
	}
}

/*
func TestWithdraw(t *testing.T) {
	// NOTICE!!
	// Substitute for your address wallet in the empty string
	success, err := c.Withdraw("DOGE", "11", "")
	if err != nil {
		t.Errorf("erro calling Withdraw endpoint: %v", err)
	}
	if success {
		t.Errorf("SUCCESS!")
	}
}

func TestWithdrawOptionalParameter(t *testing.T) {
	// NOTICE!!
	// Substitute for your address wallet in the empty string
	success, err := c.Withdraw("DOGE", "12", "", "Hello Ramon. That's my comment.")
	if err != nil {
		t.Errorf("erro calling Withdraw endpoint: %v", err)
	}
	if success {
		t.Errorf("SUCCESS!")
	}
}
*/
// func TestGetBalances_checkResults(t *testing.T) {
// 	results, _ := c.GetBalances()
// 	if len(results) == 0 {
// 		t.Errorf("could not retrieve any balance from GetBalances endpoint\n")
// 	}
// }

func TestGetMarketSummary(t *testing.T) {
	_, err := c.GetMarketSummary("ETH_BTC")
	if err != nil {
		t.Errorf("error calling getmarketsummary route: %v\n", err)
	}
}

// func TestBuyLimit(t *testing.T) {
// 	_, err := c.BuyLimit("DOGE_BTC", "100")
// 	if err != nil {
// 		t.Errorf("error trying to buy limit: %v\n", err)
// 	}
// }

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}
