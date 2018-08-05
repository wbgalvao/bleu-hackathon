package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/wbgalvao/bleu-hackathon/balance"
	"github.com/wbgalvao/bleu-hackathon/market"
	"github.com/wbgalvao/bleu-hackathon/order"
)

// Client represents a HTTP Client
type Client struct {
	BaseURL    *url.URL
	APIKey     string
	APISecret  string
	httpClient *http.Client
}

type getBalancesResponse struct {
	Success string
	Message string
	Result  []balance.Balance
}

type withdrawResponse struct {
	Success string
	Message string
	Result  []interface{}
}

type getOrdersResponse struct {
	Success string
	Message string
	Result  []order.Order
}

type marketSummaryResponse struct {
	Success string
	Message string
	Result  []market.Market
}

type limitOperationsResponse struct {
	Success string
	Message string
	Result  map[string]string
}

func getHashMacStr(message, key string) string {
	secretInBytes := []byte(key)

	mac := hmac.New(sha512.New, secretInBytes)
	mac.Write([]byte(message))
	expectedMac := mac.Sum(nil)
	fmt.Printf("[DBG] getHashMacStr %s + %s = ", message, key)

	return hex.EncodeToString(expectedMac)
}

func floatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 8, 64)
}

// BuildRequest uses the HTTP Client to build a new http.Request object
func (c *Client) BuildRequest(method, destPath string, body interface{}) (*http.Request, error) {

	u, err := url.Parse(c.BaseURL.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	u.Path = path.Join(u.Path, destPath)

	// rel := &url.URL{Path: path}
	// u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// add apiKey querystring if request is private
	q := req.URL.Query()
	q.Add("apikey", c.APIKey)

	req.URL.RawQuery = q.Encode()

	fmt.Println("[DBG] " + req.URL.String())

	return req, nil
}

// DoRequest uses the Client http.Client field to execute an http.Request
func (c *Client) DoRequest(req *http.Request, requestIsPrivate bool) (*http.Response, error) {
	if requestIsPrivate {
		q := req.URL.Query()
		str := getHashMacStr(req.URL.String(), c.APISecret)
		fmt.Println(str)
		q.Add("apisign", str)

		req.URL.RawQuery = q.Encode()
	}

	fmt.Println("[URL] " + req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// GetBalances returns a list of Balances for a given account
func (c *Client) GetBalances(opt ...string) ([]balance.Balance, error) {
	var result []balance.Balance

	if len(opt) > 1 {
		return result, fmt.Errorf("To many args for this function")
	}

	// build request
	req, err := c.BuildRequest("GET", "/account/getbalances", nil)
	if err != nil {
		return result, fmt.Errorf("error creating request for GetBalance")
	}
	fmt.Println("[DEBUG_OPT]")
	fmt.Println(opt)

	if len(opt) > 0 {
		q := req.URL.Query()
		q.Add("currencies", opt[0])
		req.URL.RawQuery = q.Encode()
	}

	// execute request
	resp, err := c.DoRequest(req, true)
	if err != nil {
		fmt.Println()
		return result, fmt.Errorf("error in GetBalances request: %v", err)
	}

	// open response body
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()

	// decode response
	var gbr getBalancesResponse
	err = json.Unmarshal(respJSON, &gbr)
	if err != nil {
		return result, fmt.Errorf("could not unmarshall response body JSON: %v", err)
	}

	fmt.Println(string(respJSON))

	if gbr.Success != "true" {
		return result, fmt.Errorf("error retrieving balance for account: %s", gbr.Message)
	}

	return gbr.Result, nil

}

func (c *Client) Withdraw(currency, quantity, destAddress string, opt ...string) (bool, error) {
	if len(opt) > 1 {
		return false, fmt.Errorf("To many args for this function")
	}
	// build request
	req, err := c.BuildRequest("GET", "/account/withdraw", nil)
	if err != nil {
		return false, fmt.Errorf("error creating request for Withdraw")
	}

	// build params for request withdraw

	q := req.URL.Query()
	q.Add("currency", currency)
	q.Add("quantity", quantity)
	q.Add("address", destAddress)

	if len(opt) > 0 {
		q.Add("comments", opt[0])
	}

	req.URL.RawQuery = q.Encode()

	// execute request
	resp, err := c.DoRequest(req, true)
	if err != nil {
		fmt.Println()
		return false, fmt.Errorf("error in Withdraw request: %v", err)
	}

	// open response body
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(string(respJSON))
	// decode response
	var wr withdrawResponse
	err = json.Unmarshal(respJSON, &wr)
	if err != nil {
		return false, fmt.Errorf("could not unmarshall response body JSON: %v", err)
	}

	if wr.Success != "true" {
		return false, fmt.Errorf("error retrieving balance for account: %s", wr.Message)
	}

	return strconv.ParseBool(wr.Success)

}

func (c *Client) ListOrder(market, orderStatus, orderType string, opt ...string) ([]order.Order, error) {
	var result []order.Order
	if len(opt) > 1 {
		return result, fmt.Errorf("To many args for this function")
	}
	// build request
	req, err := c.BuildRequest("GET", "/account/getorders", nil)
	if err != nil {
		return result, fmt.Errorf("error creating request for Withdraw")
	}

	// build params for request withdraw

	q := req.URL.Query()
	q.Add("market", market)
	q.Add("ordertype", orderType)
	q.Add("orderstatus", orderStatus)

	if len(opt) > 0 {
		q.Add("depth", opt[0])
	}

	req.URL.RawQuery = q.Encode()

	// execute request
	resp, err := c.DoRequest(req, true)
	if err != nil {
		fmt.Println()
		return result, fmt.Errorf("error in Withdraw request: %v", err)
	}

	// open response body
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(string(respJSON))
	// decode response
	var gor getOrdersResponse
	err = json.Unmarshal(respJSON, &gor)
	if err != nil {
		return result, fmt.Errorf("could not unmarshall response body JSON: %v", err)
	}

	if gor.Success != "true" {
		return result, fmt.Errorf("error retrieving balance for account: %s", gor.Message)
	}

	return gor.Result, nil
}

func (c *Client) GetMarketSummary(m string) ([]market.Market, error) {

	var result []market.Market
	if m == "" {
		return result, fmt.Errorf("empty market in function call")
	}

	// build request
	req, err := c.BuildRequest("GET", "/public/getmarketsummary", nil)
	if err != nil {
		return result, fmt.Errorf("error creating request for market summary")
	}

	// build params for request withdraw
	q := req.URL.Query()
	q.Add("market", m)
	req.URL.RawQuery = q.Encode()

	// execute request
	resp, err := c.DoRequest(req, true)
	if err != nil {
		return result, fmt.Errorf("error in market summary request: %v", err)
	}

	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()

	var msr marketSummaryResponse
	err = json.Unmarshal(respJSON, &msr)
	if err != nil {
		return result, fmt.Errorf("could not unmarshall response body JSON: %v", err)
	}

	if msr.Success != "true" {
		return result, fmt.Errorf("error retrieving market summary: %s", msr.Message)
	}

	return msr.Result, nil

}

func (c *Client) BuyLimit(m, quantity string, opt ...string) (map[string]string, error) {
	result := make(map[string]string)
	if len(opt) > 1 {
		return result, fmt.Errorf("too many args for this function")
	}

	// get market bid
	var ms []market.Market
	ms, err := c.GetMarketSummary(m)
	if err != nil {
		return result, fmt.Errorf("could not retrieve market summary: %v", err)
	}

	fmt.Println(ms)

	rate := floatToString(ms[0].Bid)

	// build request
	req, err := c.BuildRequest("GET", "/market/buylimit", nil)
	if err != nil {
		return result, fmt.Errorf("error creating request for BuyLimit")
	}

	// add querystring parameters
	q := req.URL.Query()
	q.Add("market", m)
	q.Add("rate", rate)
	q.Add("quantity", quantity)
	if len(opt) > 0 {
		q.Add("comments", opt[0])
	}
	req.URL.RawQuery = q.Encode()

	// execute request
	resp, err := c.DoRequest(req, true)
	if err != nil {
		fmt.Println()
		return result, fmt.Errorf("error in BuyLimit request: %v", err)
	}

	// open response body
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()

	// decode response
	var lor limitOperationsResponse
	err = json.Unmarshal(respJSON, &lor)
	if err != nil {
		return result, fmt.Errorf("could not unmarshall response body JSON: %v", err)
	}

	if lor.Success != "true" {
		return result, fmt.Errorf("error buying limit for account: %s", lor.Message)
	}

	return lor.Result, nil

}

func (c *Client) SellLimit(m, quantity string, opt ...string) (map[string]string, error) {
	result := make(map[string]string)
	if len(opt) > 1 {
		return result, fmt.Errorf("too many args for this function")
	}

	// get market bid
	var ms []market.Market
	ms, err := c.GetMarketSummary(m)
	if err != nil {
		return result, fmt.Errorf("could not retrieve market summary: %v", err)
	}
	rate := floatToString(ms[0].Bid)

	// build request
	req, err := c.BuildRequest("GET", "/market/selllimit", nil)
	if err != nil {
		return result, fmt.Errorf("error creating request for SellLimit")
	}

	// add querystring parameters
	q := req.URL.Query()
	q.Add("market", m)
	q.Add("rate", rate)
	q.Add("quantity", quantity)
	if len(opt) > 0 {
		q.Add("comments", opt[0])
	}
	req.URL.RawQuery = q.Encode()

	// execute request
	resp, err := c.DoRequest(req, true)
	if err != nil {
		fmt.Println()
		return result, fmt.Errorf("error in SellLimit request: %v", err)
	}

	// open response body
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(string(respJSON))

	// decode response
	var lor limitOperationsResponse
	err = json.Unmarshal(respJSON, &lor)
	if err != nil {
		return result, fmt.Errorf("could not unmarshall response body JSON: %v", err)
	}

	if lor.Success != "true" {
		return result, fmt.Errorf("error selling limit for account: %s", lor.Message)
	}

	return lor.Result, nil

}
