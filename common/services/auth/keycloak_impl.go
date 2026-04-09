package auth

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type KeycloakOAuthServiceImpl struct {
	authServerURL string
	realm         string
	jwksURL       string
	publicKeys    map[string]*rsa.PublicKey
	mu            sync.RWMutex
}

func NewKeycloakOAuthService(authServerURL, realm string) OAuthService {
	return &KeycloakOAuthServiceImpl{
		authServerURL: authServerURL,
		realm:         realm,
		jwksURL:       fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", authServerURL, realm),
		publicKeys:    make(map[string]*rsa.PublicKey),
	}
}

func (k *KeycloakOAuthServiceImpl) HandleTokenGeneration(clientId string, clientSecret string, userId string) (dto.JWTResponse, error) {
	return dto.JWTResponse{}, errors.New("token generation via Keycloak proxy not implemented; use Keycloak protocol endpoints directly")
}

func (k *KeycloakOAuthServiceImpl) ParseToken(accessToken string) (*jwt.Token, error) {
	return jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return k.getKey(token, accessToken)
	})
}

func (k *KeycloakOAuthServiceImpl) GetTokenClaims(accessToken string) (map[string]string, error) {
	token, err := k.ParseToken(accessToken)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		dataMap := make(map[string]string)
		for key, val := range claims {
			dataMap[key] = fmt.Sprintf("%v", val)
		}
		return dataMap, nil
	}
	return nil, errors.New("invalid claims")
}

func (k *KeycloakOAuthServiceImpl) IsAuthenticated(c echo.Context) (bool, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return false, nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return false, nil
	}

	return k.IsValidToken(parts[1])
}

func (k *KeycloakOAuthServiceImpl) IsValidToken(accessToken string) (bool, error) {
	token, err := k.ParseToken(accessToken)
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return false, err
	}
	return token.Valid, nil
}

func (k *KeycloakOAuthServiceImpl) getKey(token *jwt.Token, accessToken string) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid header not found")
	}

	k.mu.RLock()
	key, exists := k.publicKeys[kid]
	k.mu.RUnlock()

	if exists {
		return key, nil
	}

	// Fetch JWKS if key not found
	if err := k.refreshKeys(accessToken); err != nil {
		return nil, err
	}

	k.mu.RLock()
	key, exists = k.publicKeys[kid]
	k.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("key not found for kid: %s", kid)
	}

	return key, nil
}

func (k *KeycloakOAuthServiceImpl) refreshKeys(requestToken string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	log.Printf("Refreshing JWKS from %s", k.jwksURL)

	req, err := http.NewRequest("GET", k.jwksURL, nil)
	if err != nil {
		return err
	}

	if requestToken != "" {
		req.Header.Set("Authorization", "Bearer "+requestToken)
	} else {
		log.Printf("Warning: No token available for JWKS refresh: %v. Attempting unauthenticated request.", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kid string   `json:"kid"`
			Kty string   `json:"kty"`
			Alg string   `json:"alg"`
			Use string   `json:"use"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			X5c []string `json:"x5c"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	for _, keyInfo := range jwks.Keys {
		if len(keyInfo.X5c) > 0 {
			cert := "-----BEGIN CERTIFICATE-----\n" + keyInfo.X5c[0] + "\n-----END CERTIFICATE-----"
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			if err == nil {
				k.publicKeys[keyInfo.Kid] = publicKey
			}
		}
	}

	return nil
}
