package impl

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/MarkoLuna/AuthService/internal/models"
	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type clientServiceStub struct {
	validClientId     string
	validClientSecret string
}

func (s *clientServiceStub) IsValidClientCredentials(clientId string, clientSecret string) (bool, error) {
	if clientId == s.validClientId && clientSecret == s.validClientSecret {
		return true, nil
	}
	return false, nil
}

type oauthServiceStub struct {
	response dto.JWTResponse
	err      error
}

func (s *oauthServiceStub) HandleTokenGeneration(clientId string, clientSecret string, userId string) (dto.JWTResponse, error) {
	return s.response, s.err
}

func (s *oauthServiceStub) ParseToken(accessToken string) (*jwt.Token, error) {
	return nil, nil
}

func (s *oauthServiceStub) GetTokenClaims(accessToken string) (map[string]string, error) {
	return map[string]string{
		"sub":  "user-1",
		"role": "admin",
	}, nil
}

func (s *oauthServiceStub) IsAuthenticated(c echo.Context) (bool, error) {
	return true, nil
}

func (s *oauthServiceStub) IsValidToken(accessToken string) (bool, error) {
	return true, nil
}

func TestLocalAuthStrategy_HandleTokenGeneration_Success(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	userSvc := NewUserService(repo)
	clientSvc := &clientServiceStub{
		validClientId:     "my-client",
		validClientSecret: "my-secret",
	}
	oauthSvc := &oauthServiceStub{
		response: dto.JWTResponse{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresIn:    120,
			Scope:        "all",
			TokenType:    "Bearer",
		},
	}

	preloadUser(repo, models.User{
		Id:           "user-1",
		Username:     "testuser",
		PasswordHash: hashPassword(t, "pass123"),
		Enabled:      true,
	})

	strategy := NewLocalAuthStrategy(clientSvc, userSvc, oauthSvc)

	e := echo.New()
	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("password", "pass123")

	auth := base64.StdEncoding.EncodeToString([]byte("my-client:my-secret"))
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+auth)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	resp, err := strategy.HandleTokenGeneration(c)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-access-token", resp.AccessToken)
	assert.Equal(t, "Bearer", resp.TokenType)
}

func TestLocalAuthStrategy_HandleTokenGeneration_MissingAuthHeader(t *testing.T) {
	svc := NewLocalAuthStrategy(nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	resp, err := svc.HandleTokenGeneration(c)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "missing authorization header")
}

func TestLocalAuthStrategy_HandleTokenGeneration_InvalidAuthType(t *testing.T) {
	svc := NewLocalAuthStrategy(nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	resp, err := svc.HandleTokenGeneration(c)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid authorization type")
}

func TestLocalAuthStrategy_HandleTokenGeneration_InvalidBase64(t *testing.T) {
	svc := NewLocalAuthStrategy(nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", nil)
	req.Header.Set("Authorization", "Basic not-valid-base64!!!")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	resp, err := svc.HandleTokenGeneration(c)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid authorization header")
}

func TestLocalAuthStrategy_HandleTokenGeneration_InvalidClient(t *testing.T) {
	clientSvc := &clientServiceStub{
		validClientId:     "real-client",
		validClientSecret: "real-secret",
	}
	svc := NewLocalAuthStrategy(clientSvc, nil, nil)

	e := echo.New()
	auth := base64.StdEncoding.EncodeToString([]byte("wrong-client:wrong-secret"))
	req := httptest.NewRequest(http.MethodPost, "/oauth/token", nil)
	req.Header.Set("Authorization", "Basic "+auth)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	resp, _ := svc.HandleTokenGeneration(c)
	assert.Nil(t, resp)
}

func TestLocalAuthStrategy_GetTokenClaims_Success(t *testing.T) {
	oauthSvc := &oauthServiceStub{}
	strategy := NewLocalAuthStrategy(nil, nil, oauthSvc)

	claims, err := strategy.GetTokenClaims("some-token")
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-1", claims["sub"])
}
