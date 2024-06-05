package config

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type Config struct {
	PassKey        string
	ConsumerKey    string
	ConsumerSecret string
}

func NewConfig(passKey, consumerKey, consumerSecret string) *Config {
	return &Config{
		PassKey:        passKey,
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}
}

func (c *Config) getPassKey() string {
	return c.PassKey
}

func (c *Config) getAuth() (string, error) {
	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	authorizationString := base64.StdEncoding.EncodeToString([]byte(c.ConsumerKey + ":" + c.ConsumerSecret))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+authorizationString)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var accessTokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}
	if err := json.NewDecoder(res.Body).Decode(&accessTokenResponse); err != nil {
		return "", err
	}

	return accessTokenResponse.AccessToken, nil
}