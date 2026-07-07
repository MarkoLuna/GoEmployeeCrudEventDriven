package config

import (
	"net/http"
	"strings"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"
)

var (
	DefaultSkippedPaths = [...]string{
		"/healthcheck/",
		"/swagger/",
	}
)

type AuthConfig struct {
	EnableAuth       bool
	SkippedPaths     []string
	ValidationClient *auth.TokenValidationClient
}

func NewAuthConfig(echoInstance *echo.Echo, enableAuth bool, skippedPaths []string, validationClient *auth.TokenValidationClient) {
	if enableAuth {
		authConfig := AuthConfig{
			EnableAuth:       enableAuth,
			SkippedPaths:     skippedPaths,
			ValidationClient: validationClient,
		}

		echoInstance.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if authConfig.isSkippedPath(c.Request().URL.Path) {
					return next(c)
				}

				accessToken, ok := utils.GetBearerAuth(c.Request().Header)
				if !ok {
					return c.String(http.StatusUnauthorized, "invalid token")
				}

				claims, err := authConfig.ValidationClient.ValidateToken(accessToken)
				if err != nil {
					return c.String(http.StatusUnauthorized, err.Error())
				}

				c.Set("claims", claims)
				return next(c)
			}
		})
	}
}

func (authConfig AuthConfig) isSkippedPath(path string) bool {
	for i := 0; i < len(authConfig.SkippedPaths); i++ {
		if strings.HasPrefix(path, authConfig.SkippedPaths[i]) {
			return true
		}
	}
	return false
}
