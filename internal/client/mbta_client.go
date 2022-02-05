package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	MbtaApiBaseUrl       = "https://api-v3.mbta.com"
	MbtaApiBaseUrlFormat = MbtaApiBaseUrl + "%s"
	MbtaApiKeyHeader     = "x-api-key"
	// MbtaApiKey would usually be stored in a secrets manager, but for the purposes of this exercise.
	MbtaApiKey = "e416deedab3641faa20128f300305998"

	StationLocationType = 1
)

type MbtaClient struct {
	client http.Client
}

func NewMbtaApiClient() *MbtaClient {
	return &MbtaClient{
		http.Client{Timeout: time.Duration(20) * time.Second},
	}
}

func (c *MbtaClient) GetRoutes() (*RoutesResponse, error) {
	req, err := c.newRequest("GET", "/routes")
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("filter[type]", "0,1")
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Request failed (%s): %s", resp.Status, string(body)),
		)
	}

	routes := &RoutesResponse{}
	err = json.Unmarshal(body, routes)
	return routes, err
}

func (c *MbtaClient) GetStopsForRoutes(routeId string) (*StopsResponse, error) {
	req, err := c.newRequest("GET", "/stops")
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	// Only use stations for this fetch
	q.Add("filter[route]", routeId)
	// Only fetch stations for this route
	q.Add("filter[location_type]", strconv.FormatInt(StationLocationType, 10))
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Request failed (%s): %s", resp.Status, string(body)),
		)
	}

	stops := &StopsResponse{}
	err = json.Unmarshal(body, stops)
	return stops, err
}

func (c *MbtaClient) newRequest(method string, endpoint string) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf(MbtaApiBaseUrlFormat, endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Header[MbtaApiKeyHeader] = []string{MbtaApiKey}
	return req, nil
}
