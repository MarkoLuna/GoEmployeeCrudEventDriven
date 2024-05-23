package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheckHandler Healthcheck
// @Tags 	Healthcheck
// @Summary      healthcheck
// @Description  get healthcheck status
// @Accept       json
// @Produce      json
// @Success      200  {string}  OK
// @Router       /healthcheck/ [get]
func HealthCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
