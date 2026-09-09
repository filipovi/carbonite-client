package api

import (
	"carbonite/admin/internal/data"
	"encoding/json"
	"strings"
)

func (client Client) PostToken(email string, password string) (*data.TokenPayload, error) {
	body, err := json.Marshal(data.TokenBody{
		Email:        email,
		Password:     password,
		ClientID:     client.Config.ClientID,
		ClientSecret: client.Config.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	d, err := client.handleRequest(strings.Join([]string{client.Config.URL, "/v1/authentication"}, ""), "POST", body, nil)
	if err != nil {
		return nil, err
	}

	var payload data.TokenPayload

	json.Unmarshal(d, &payload)
	return &payload, nil
}

func (client Client) PostRefreshToken(refresh string) (*data.TokenPayload, error) {
	body, err := json.Marshal(data.RefreshTokenBody{
		Refresh: refresh,
	})
	if err != nil {
		return nil, err
	}

	d, err := client.handleRequest(strings.Join([]string{client.Config.URL, "/v1/authentication/refresh"}, ""), "POST", body, nil)
	if err != nil {
		return nil, err
	}

	var payload data.TokenPayload

	json.Unmarshal(d, &payload)
	return &payload, nil
}
