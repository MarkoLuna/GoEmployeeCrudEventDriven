package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TokenValidationClient struct {
	authServiceURL string
	httpClient     *http.Client
}

func NewTokenValidationClient(authServiceURL string) *TokenValidationClient {
	return &TokenValidationClient{
		authServiceURL: authServiceURL,
		httpClient:     &http.Client{},
	}
}

func (c *TokenValidationClient) ValidateToken(accessToken string) (map[string]string, error) {
	req, err := http.NewRequest("GET", c.authServiceURL+"/oauth/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("authentication service unavailable")
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("authentication service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("authentication service unavailable")
	}

	var claims map[string]string
	if err := json.Unmarshal(body, &claims); err != nil {
		return nil, fmt.Errorf("authentication service unavailable")
	}

	return claims, nil
}
