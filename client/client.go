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

	"github.com/wbgalvao/bleu-hackathon/balance"
)

// Client represents a HTTP Client
type Client struct {
	BaseURL    *url.URL
	APIKey     string
	APISecret  string
	httpClient *http.Client
}

type getBalancesResponse struct {
	success bool
	message string
	result  []balance.Balance
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

	return req, nil
}

// DoRequest uses the Client http.Client field to execute an http.Request
func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func getMacStr(message, key string) string {
	secretInBytes := []byte(c.APISecret)

	mac := hmac.New(sha512.New, secretInBytes)
	mac.Write([]byte(message))
	expectedMac := mac.Sum(nil)
	fmt.Printf("[DBG] %s + %s = ", message, key)

	return hex.EncodeToString(expectedMac)
}

// GetBalances returns a list of Balances for a given account
func (c *Client) GetBalances() ([]balance.Balance, error) {
	var result []balance.Balance

	// build request
	req, err := c.BuildRequest("GET", "/account/getbalances", nil)
	if err != nil {
		return result, fmt.Errorf("error creating request for GetBalance")
	}

	// add apiKey querystring
	q := req.URL.Query()
	q.Add("apikey", c.APIKey)

	req.URL.RawQuery = q.Encode()

	str := getMacStr(req.URL.String(), c.APISecret)
	fmt.Println(str)
	q.Add("apisign", str)

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	// execute request
	resp, err := c.DoRequest(req)
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

	if gbr.success != true {
		return result, fmt.Errorf("error retrieving balance for account: %v", gbr.message)
	}

	return gbr.result, nil

}
