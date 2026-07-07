package impl

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/MarkoLuna/AuthService/internal/services"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/labstack/echo/v4"
)

type LocalAuthStrategy struct {
	clientService services.ClientService
	userService   services.UserService
	oauthService  services.OAuthService
}

func NewLocalAuthStrategy(
	clientService services.ClientService,
	userService services.UserService,
	oauthService services.OAuthService,
) *LocalAuthStrategy {
	return &LocalAuthStrategy{
		clientService: clientService,
		userService:   userService,
		oauthService:  oauthService,
	}
}

func (s *LocalAuthStrategy) HandleTokenGeneration(c echo.Context) (*dto.JWTResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		return nil, errors.New("invalid authorization type")
	}

	authDecoded, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
	if err != nil {
		return nil, errors.New("invalid authorization header")
	}

	parts := strings.SplitN(string(authDecoded), ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid authorization format")
	}
	clientIdReq := parts[0]
	clientSecretReq := parts[1]

	validClientCred, err := s.clientService.IsValidClientCredentials(clientIdReq, clientSecretReq)
	if !validClientCred || err != nil {
		return nil, err
	}

	userNameReq := c.FormValue("username")
	passwordReq := c.FormValue("password")

	userId, err := s.userService.GetUserId(userNameReq, passwordReq)
	if err != nil {
		return nil, err
	}

	jWTResponse, err := s.oauthService.HandleTokenGeneration(clientIdReq, clientSecretReq, userId)
	if err != nil {
		return nil, err
	}

	return &jWTResponse, nil
}

func (s *LocalAuthStrategy) GetTokenClaims(accessToken string) (map[string]string, error) {
	return s.oauthService.GetTokenClaims(accessToken)
}
