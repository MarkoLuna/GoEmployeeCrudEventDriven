package services

import (
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/labstack/echo/v4"
)

type AuthStrategy interface {
	HandleTokenGeneration(c echo.Context) (*dto.JWTResponse, error)
	GetTokenClaims(accessToken string) (map[string]string, error)
}
