package config

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func EnableCORS(echoInstance *echo.Echo) {
	echoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderAccept,
			echo.HeaderContentType,
			echo.HeaderContentLength,
			echo.HeaderAcceptEncoding,
			echo.HeaderAuthorization,
			echo.HeaderXCSRFToken,
		},
		AllowCredentials: true,
		AllowMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodOptions,
			http.MethodPut,
			http.MethodDelete,
		},
	}))
}
