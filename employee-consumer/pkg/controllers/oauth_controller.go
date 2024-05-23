package controllers

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/services"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/server"

	"github.com/labstack/echo/v4"
)

type OAuthController struct {
	oauthServer   *server.Server
	oauthSevice   services.OAuthService
	clientService services.ClientService
	userService   services.UserService
}

func NewOAuthController(oauthServer *server.Server,
	oauthService services.OAuthService,
	clientService services.ClientService,
	userService services.UserService) OAuthController {
	ctrl := OAuthController{oauthServer, oauthService, clientService, userService}
	ctrl.Configure()
	return ctrl
}

func (ctrl OAuthController) Configure() {

	ctrl.oauthServer.SetAllowGetAccessRequest(true)
	ctrl.oauthServer.SetClientInfoHandler(server.ClientFormHandler)

	ctrl.oauthServer.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	ctrl.oauthServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		re.SetHeader("error", err.Error())
		return
	})
}

// Login auth
// @Tags 	Auth
// @Summary login user
// @Description login user
// @Accept  json
// @Produce  json
// @Param   password      query string true  "Password"
// @Param   username      query string true  "Username"
// @Param   grant_type      path string true  "Grant type"
// @Success 200 {object} dto.JWTResponse	"ok"
// @Failure 400 {object} error "Invalid authorization!!"
// @Security BasicAuth
// @Router /oauth/token [post]
func (ctrl OAuthController) TokenHandler(c echo.Context) error {
	auth, ok := utils.GetBasicAuth(c.Request().Header)
	if !ok {
		return c.String(http.StatusUnauthorized, "Unable to find the Authentication")
	}

	clientIdReq, clientSecretReq := ctrl.DecodeBasicAuth(auth)
	validClientCred, err := ctrl.clientService.IsValidClientCredentials(clientIdReq, clientSecretReq)
	if !validClientCred || err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	userNameReq := c.FormValue("username")
	passwordReq := c.FormValue("password")

	userId, err := ctrl.userService.GetUserId(userNameReq, passwordReq)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	jWTResponse, err := ctrl.oauthSevice.HandleTokenGeneration(clientIdReq, clientSecretReq, userId)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, jWTResponse)
}

// GetUserInfo auth
// @Tags 	Auth
// @Summary get-user-info
// @Description get user info
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string	"ok"
// @Failure 400 {object} error "Invalid authorization!!"
// @Security ApiKeyAuth
// @Router /oauth/userinfo [get]
func (ctrl OAuthController) GetUserInfo(c echo.Context) error {
	accessToken, ok := utils.GetBearerAuth(c.Request().Header)
	if !ok {
		return c.String(http.StatusUnauthorized, "Unable to find the Authentication")
	}

	log.Println("auth token: ", accessToken)
	claims, err := ctrl.oauthSevice.GetTokenClaims(accessToken)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, claims)
}

func (ctrl OAuthController) DecodeBasicAuth(auth string) (string, string) {
	authDecoded, _ := base64.StdEncoding.DecodeString(auth)
	authReq := strings.Split(string(authDecoded), ":")

	return authReq[0], authReq[1]
}
