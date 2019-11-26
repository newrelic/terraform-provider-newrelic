package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// NewQueryClient makes a new client for the user to query with.
func NewQueryClient(queryKey, accountID string) *QueryClient {
	client := &QueryClient{}
	client.URL = createQueryURL(accountID)
	client.QueryKey = queryKey
	client.Logger = log.New()

	// Defaults
	client.RequestTimeout = DefaultQueryRequestTimeout
	client.RetryCount = DefaultRetries
	client.RetryWait = DefaultRetryWaitTime

	return client
}

func createQueryURL(accountID string) *url.URL {
	insightsURL, _ := url.Parse(insightsQueryURL)
	insightsURL.Path = fmt.Sprintf("%s/%s/query", insightsURL.Path, accountID)
	return insightsURL
}

// Validate makes sure the QueryClient is configured correctly for use
func (c *QueryClient) Validate() error {
	if correct, _ := regexp.MatchString("api.newrelic.com/v1/accounts/[0-9]+/query", c.URL.String()); !correct {
		return fmt.Errorf("invalid query endpoint %s", c.URL)
	}

	if len(c.QueryKey) < 1 {
		return fmt.Errorf("not a valid license key: %s", c.QueryKey)
	}
	return nil
}

// QueryEvents initiates an Insights query, returns a response for parsing
func (c *QueryClient) QueryEvents(nrqlQuery string) (response *QueryResponse, err error) {
	response = &QueryResponse{}
	err = c.Query(nrqlQuery, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Query initiates an Insights query, with the JSON parsed into 'response' struct
func (c *QueryClient) Query(nrqlQuery string, response interface{}) (err error) {
	if response == nil {
		return errors.New("go-insights: Invalid query response can not be nil")
	}

	err = c.queryRequest(nrqlQuery, response)
	if err != nil {
		return err
	}
	return nil
}

// queryRequest makes a NRQL query and returns the result in `queryResult`
// which must be a pointer to a struct that the JSON package can unmarshall
func (c *QueryClient) queryRequest(nrqlQuery string, queryResult interface{}) (err error) {
	var request *http.Request
	var response *http.Response

	queryURL, err := c.generateQueryURL(nrqlQuery)
	if err != nil {
		return err
	}

	if queryResult == nil {
		return errors.New("must have pointer for result")
	}

	request, err = http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Query-Key", c.QueryKey)

	client := &http.Client{Timeout: c.RequestTimeout}

	response, err = client.Do(request)
	if err != nil {
		err = fmt.Errorf("failed query request for: %v", err)
		return
	}
	defer func() {
		respErr := response.Body.Close()
		if respErr != nil && err == nil {
			err = respErr // Don't mask previous errors
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response code: %d", response.StatusCode)
		return
	}

	err = c.parseResponse(response, queryResult)
	if err != nil {
		err = fmt.Errorf("failed query: %v", err)
	}

	return err
}

// generateQueryURL URL encodes the NRQL
func (c *QueryClient) generateQueryURL(nrqlQuery string) (string, error) {
	if len(nrqlQuery) < minValidNRQLLength {
		fmt.Println("Query was too short")
		return "", fmt.Errorf("NRQL query is too short [%s]", nrqlQuery)
	}

	// Use a new set of Values to sanitize the query string
	urlQuery := url.Values{}
	urlQuery.Set("nrql", nrqlQuery)
	queryString := urlQuery.Encode()

	queryURL := c.URL.String() + "?" + queryString

	c.Logger.Debugf("query url is: %s", queryURL)

	return queryURL, nil
}

// parseQueryResponse takes an HTTP response, make sure it is a valid response,
// then attempts to decode the JSON body into the `parsedResponse` interface
func (c *QueryClient) parseResponse(response *http.Response, parsedResponse interface{}) error {
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read response body: %s", readErr.Error())
	}

	c.Logger.Debugf("Response %d body: %s", response.StatusCode, body)

	if jsonErr := json.Unmarshal(body, parsedResponse); jsonErr != nil {
		return fmt.Errorf("unable to unmarshal query response: %v", jsonErr)
	}

	return nil
}
