package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TokenValidationClient struct {
	authServiceURL string
	httpClient     *http.Client
}

func NewTokenValidationClient(authServiceURL string, timeout time.Duration) *TokenValidationClient {
	return &TokenValidationClient{
		authServiceURL: authServiceURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *TokenValidationClient) ValidateToken(accessToken string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.httpClient.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.authServiceURL+"/oauth/userinfo", nil)
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
