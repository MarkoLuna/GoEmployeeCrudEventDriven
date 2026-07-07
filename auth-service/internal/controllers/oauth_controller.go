package controllers

import (
	"net/http"

	"github.com/MarkoLuna/AuthService/internal/services"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"
)

type OAuthController struct {
	authStrategy services.AuthStrategy
}

func NewOAuthController(authStrategy services.AuthStrategy) OAuthController {
	return OAuthController{authStrategy}
}

// Login auth
// @Tags Auth
// @Summary login user
// @Description login user
// @Accept json
// @Produce json
// @Param password formData string true "Password"
// @Param username formData string true "Username"
// @Success 200 {object} dto.JWTResponse "ok"
// @Failure 400 {object} error "Invalid authorization!!"
// @Security BasicAuth
// @Router /oauth/token [post]
func (ctrl OAuthController) TokenHandler(c echo.Context) error {
	response, err := ctrl.authStrategy.HandleTokenGeneration(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

// GetUserInfo auth
// @Tags Auth
// @Summary get-user-info
// @Description get user info
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "ok"
// @Failure 400 {object} error "Invalid authorization!!"
// @Security ApiKeyAuth
// @Router /oauth/userinfo [get]
func (ctrl OAuthController) GetUserInfo(c echo.Context) error {
	accessToken, ok := utils.GetBearerAuth(c.Request().Header)
	if !ok {
		return c.String(http.StatusUnauthorized, "Unable to find the Authentication")
	}

	claims, err := ctrl.authStrategy.GetTokenClaims(accessToken)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, claims)
}
