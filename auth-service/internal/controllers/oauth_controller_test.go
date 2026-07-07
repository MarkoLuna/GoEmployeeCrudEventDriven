package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type authStrategySuccessStub struct{}

func (s *authStrategySuccessStub) HandleTokenGeneration(c echo.Context) (*dto.JWTResponse, error) {
	return &dto.JWTResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    120,
		Scope:        "all",
		TokenType:    "Bearer",
	}, nil
}

func (s *authStrategySuccessStub) GetTokenClaims(accessToken string) (map[string]string, error) {
	return map[string]string{
		"sub": "user-1", "role": "admin",
	}, nil
}

type authStrategyUnauthorizedStub struct{}

func (s *authStrategyUnauthorizedStub) HandleTokenGeneration(c echo.Context) (*dto.JWTResponse, error) {
	return nil, errors.New("invalid credentials")
}

func (s *authStrategyUnauthorizedStub) GetTokenClaims(accessToken string) (map[string]string, error) {
	return nil, errors.New("invalid token")
}

func TestOAuthController_TokenHandler_Success(t *testing.T) {
	strategy := &authStrategySuccessStub{}
	ctrl := NewOAuthController(strategy)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.TokenHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp dto.JWTResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "access-token", resp.AccessToken)
		assert.Equal(t, "Bearer", resp.TokenType)
	}
}

func TestOAuthController_TokenHandler_Unauthorized(t *testing.T) {
	strategy := &authStrategyUnauthorizedStub{}
	ctrl := NewOAuthController(strategy)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.TokenHandler(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}

func TestOAuthController_GetUserInfo_Success(t *testing.T) {
	strategy := &authStrategySuccessStub{}
	ctrl := NewOAuthController(strategy)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.GetUserInfo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var claims map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &claims)
		assert.NoError(t, err)
		assert.Equal(t, "user-1", claims["sub"])
	}
}

func TestOAuthController_GetUserInfo_NoAuthHeader(t *testing.T) {
	strategy := &authStrategySuccessStub{}
	ctrl := NewOAuthController(strategy)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oauth/userinfo", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.GetUserInfo(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}

func TestOAuthController_GetUserInfo_InvalidToken(t *testing.T) {
	strategy := &authStrategyUnauthorizedStub{}
	ctrl := NewOAuthController(strategy)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.GetUserInfo(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}
