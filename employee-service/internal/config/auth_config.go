package config

import (
	"errors"
	"strings"

	"github.com/MarkoLuna/EmployeeService/internal/services"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"

	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
)

var (
	DefaultSkippedPaths = [...]string{
		"/healthcheck/",
		"/oauth/token",
		"/swagger/",
	}
)

var (
	signingKey = utils.GetEnv("OAUTH_SIGNING_KEY", "00000000")
)

type AuthConfig struct {
	EnableAuth   bool
	SkippedPaths []string
	OAuthService services.OAuthService
}

func NewAuthConfig(echoInstance *echo.Echo, enableAuth bool, skippedPaths []string, authService services.OAuthService) {
	if enableAuth {
		authConfig := AuthConfig{EnableAuth: enableAuth, SkippedPaths: skippedPaths, OAuthService: authService}

		defaultJWTConfig := echojwt.Config{
			SigningKey: []byte(signingKey),
			Skipper: func(e echo.Context) bool {
				return authConfig.isSkippedPath(e.Request().URL.Path)
			},
			ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {

				accessToken, ok := utils.GetBearerAuth(c.Request().Header)
				if !ok {
					return nil, errors.New("invalid token")
				}
				token, err := authConfig.OAuthService.ParseToken(accessToken)
				if err != nil {
					return nil, err
				}

				if !token.Valid {
					return nil, errors.New("invalid token")
				}
				return token, nil
			},
		}

		echoInstance.Use(echojwt.WithConfig(defaultJWTConfig))
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
