package services

import (
	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
)

type OAuthService interface {
	HandleTokenGeneration(clientId string, clientSecret string, userId string) (dto.JWTResponse, error)
	ParseToken(accessToken string) (*jwt.Token, error)
	GetTokenClaims(accessToken string) (map[string]string, error)
	IsAuthenticated(c echo.Context) (bool, error)
	IsValidToken(accessToken string) (bool, error)
}
