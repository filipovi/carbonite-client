package api

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"slices"

	"carbonite/client/internal/data"
)

type (
	options struct{}

	Option func(options *options) error

	cfg struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		URL          string `json:"URL"`
	}

	// Client is the API client structure
	Client struct {
		Client *http.Client
		Config cfg
	}

	// API interface
	API any
)

// New creates a news API Client
func New(clientID string, clientSecret string, apiURL string, env string) (*Client, error) {
	var tr *http.Transport

	if env == "prod" {
		tr = &http.Transport{}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &Client{
		Config: cfg{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			URL:          apiURL,
		},
		Client: &http.Client{Transport: tr},
	}, nil
}

func (client Client) handleRequest(url string, verb string, b []byte, token *data.TokenPayload) ([]byte, error) {
	requestBody := bytes.NewBuffer(b)
	req, err := http.NewRequest(verb, url, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AuthenticationToken.Token))
	}
	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	errorCodes := []int{200, 201, 206}
	if !slices.Contains(errorCodes, resp.StatusCode) {
		return nil, fmt.Errorf("error code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
