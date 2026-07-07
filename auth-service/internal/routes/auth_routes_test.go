package routes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MarkoLuna/AuthService/internal/controllers"
	"github.com/MarkoLuna/AuthService/internal/services/stubs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAuthRoutes(t *testing.T) {
	echoInstance := echo.New()

	authStrategy := stubs.NewAuthStrategyStub()
	oauthController := controllers.NewOAuthController(authStrategy)

	RegisterAuthRoutes(echoInstance, &oauthController)

	tables := []struct {
		method      string
		path        string
		status      int
		handlerName string
	}{
		{"POST", "/oauth/token", http.StatusOK, "TokenHandler"},
		{"GET", "/oauth/userinfo", http.StatusUnauthorized, "GetUserInfo"},
	}

	for _, table := range tables {
		req, err := http.NewRequest(table.method, table.path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := echoInstance.NewContext(req, rec)

		switch table.handlerName {
		case "TokenHandler":
			if assert.NoError(t, oauthController.TokenHandler(c)) {
				assert.Equal(t, table.status, rec.Code, "handler returned wrong status code for %s %s", table.method, table.path)
			}
		case "GetUserInfo":
			if assert.NoError(t, oauthController.GetUserInfo(c)) {
				assert.Equal(t, table.status, rec.Code, "handler returned wrong status code for %s %s", table.method, table.path)
			}
		}
	}
}

func TestRegisterUserRoutes(t *testing.T) {
	echoInstance := echo.New()

	userService := stubs.NewUserServiceStub()
	userController := controllers.NewUserController(userService)

	RegisterUserRoutes(echoInstance, &userController)

	tables := []struct {
		method      string
		path        string
		body        string
		status      int
		handlerName string
	}{
		{"POST", "/api/user/", `{"username":"u","password":"p"}`, http.StatusCreated, "CreateUser"},
		{"GET", "/api/user/", "", http.StatusOK, "GetUsers"},
		{"GET", "/api/user/1", "", http.StatusNotFound, "GetUserById"},
		{"PUT", "/api/user/1", `{"username":"u","password":"p"}`, http.StatusNotFound, "UpdateUser"},
		{"DELETE", "/api/user/1", "", http.StatusNoContent, "DeleteUser"},
	}

	for _, table := range tables {
		var bodyReader io.Reader
		if table.body != "" {
			bodyReader = strings.NewReader(table.body)
		}
		req, err := http.NewRequest(table.method, table.path, bodyReader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := echoInstance.NewContext(req, rec)

		if table.handlerName == "GetUserById" || table.handlerName == "UpdateUser" || table.handlerName == "DeleteUser" {
			c.SetParamNames("userId")
			c.SetParamValues("1")
		}

		switch table.handlerName {
		case "CreateUser":
			if assert.NoError(t, userController.CreateUser(c)) {
				assert.Equal(t, table.status, rec.Code)
			}
		case "GetUsers":
			if assert.NoError(t, userController.GetUsers(c)) {
				assert.Equal(t, table.status, rec.Code)
			}
		case "GetUserById":
			if assert.NoError(t, userController.GetUserById(c)) {
				assert.Equal(t, table.status, rec.Code)
			}
		case "UpdateUser":
			if assert.NoError(t, userController.UpdateUser(c)) {
				assert.Equal(t, table.status, rec.Code)
			}
		case "DeleteUser":
			if assert.NoError(t, userController.DeleteUser(c)) {
				assert.Equal(t, table.status, rec.Code)
			}
		}
	}
}

func TestRegisterHealthcheckRoute(t *testing.T) {
	echoInstance := echo.New()
	RegisterHealthcheckRoute(echoInstance)

	req := httptest.NewRequest(http.MethodGet, "/healthcheck/", nil)
	rec := httptest.NewRecorder()
	echoInstance.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Healthy", rec.Body.String())
}
