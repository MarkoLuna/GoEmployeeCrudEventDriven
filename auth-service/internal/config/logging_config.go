package config

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigureLogging() {
}

func ConfigureLoggingMiddleware(echoInstance *echo.Echo) {
	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
}
