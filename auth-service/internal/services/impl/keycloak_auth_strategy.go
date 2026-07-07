package impl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	commonAuth "github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/labstack/echo/v4"
)

type KeycloakAuthStrategy struct {
	authServerURL string
	realm         string
	kcOAuthService commonAuth.OAuthService
}

func NewKeycloakAuthStrategy(authServerURL string, realm string) *KeycloakAuthStrategy {
	return &KeycloakAuthStrategy{
		authServerURL:  authServerURL,
		realm:          realm,
		kcOAuthService: commonAuth.NewKeycloakOAuthService(authServerURL, realm),
	}
}

func (s *KeycloakAuthStrategy) HandleTokenGeneration(c echo.Context) (*dto.JWTResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		return nil, errors.New("invalid authorization type")
	}

	authDecoded, err := decodeBasicAuth(authHeader)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(authDecoded, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid authorization format")
	}
	clientId := parts[0]
	clientSecret := parts[1]

	username := c.FormValue("username")
	password := c.FormValue("password")
	grantType := c.FormValue("grant_type")
	if grantType == "" {
		grantType = "password"
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.authServerURL, s.realm)

	data := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"grant_type":    {grantType},
		"username":      {username},
		"password":      {password},
	}

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("keycloak token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read keycloak response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak returned error: %s", string(body))
	}

	var kcResponse struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		RefreshToken     string `json:"refresh_token"`
		TokenType        string `json:"token_type"`
		Scope            string `json:"scope"`
	}

	if err := json.Unmarshal(body, &kcResponse); err != nil {
		return nil, fmt.Errorf("failed to parse keycloak response: %w", err)
	}

	return &dto.JWTResponse{
		AccessToken:  kcResponse.AccessToken,
		RefreshToken: kcResponse.RefreshToken,
		ExpiresIn:    int64(kcResponse.ExpiresIn),
		Scope:        kcResponse.Scope,
		TokenType:    kcResponse.TokenType,
	}, nil
}

func (s *KeycloakAuthStrategy) GetTokenClaims(accessToken string) (map[string]string, error) {
	return s.kcOAuthService.GetTokenClaims(accessToken)
}

func decodeBasicAuth(authHeader string) (string, error) {
	encoded := authHeader[len("Basic "):]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("invalid authorization header")
	}
	return string(decoded), nil
}
