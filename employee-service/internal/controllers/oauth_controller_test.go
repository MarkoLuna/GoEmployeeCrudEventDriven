package controllers

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/MarkoLuna/EmployeeService/internal/services"
	"github.com/MarkoLuna/EmployeeService/internal/services/stubs"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupOAuthController() OAuthController {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	oauthServer := server.NewDefaultServer(manager)
	
	oauthService := stubs.NewOAuthServiceStub()
	clientService := services.NewClientService()
	userService := services.NewUserService()
	
	return NewOAuthController(oauthServer, oauthService, clientService, userService)
}

func TestTokenHandler_Success(t *testing.T) {
	e := echo.New()
	
	f := make(url.Values)
	f.Set("username", "user")
	f.Set("password", "secret")
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	basicAuth := base64.StdEncoding.EncodeToString([]byte("client:password"))
	req.Header.Set("Authorization", "Basic "+basicAuth)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctrl := setupOAuthController()

	err := ctrl.TokenHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTokenHandler_InvalidBasicAuth(t *testing.T) {
	e := echo.New()
	
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", nil)
	// Missing basic auth header

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctrl := setupOAuthController()

	ctrl.TokenHandler(c)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetUserInfo_Success(t *testing.T) {
	e := echo.New()
	
	req := httptest.NewRequest(http.MethodGet, "/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctrl := setupOAuthController()

	err := ctrl.GetUserInfo(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserInfo_NoToken(t *testing.T) {
	e := echo.New()
	
	req := httptest.NewRequest(http.MethodGet, "/oauth/userinfo", nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctrl := setupOAuthController()

	ctrl.GetUserInfo(c)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
