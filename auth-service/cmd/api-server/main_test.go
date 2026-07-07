package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/MarkoLuna/AuthService/internal/controllers"
	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/MarkoLuna/AuthService/internal/routes"
	"github.com/MarkoLuna/AuthService/internal/services/impl"
	"github.com/MarkoLuna/AuthService/internal/services/stubs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var basePath string

func InitServer() string {
	App.EchoInstance = echo.New()
	App.UserRepository = repositories.NewUserRepositoryStub()
	App.UserService = impl.NewUserService(App.UserRepository)
	App.AuthStrategy = stubs.NewAuthStrategyStub()
	App.OAuthController = controllers.NewOAuthController(App.AuthStrategy)
	App.UserController = controllers.NewUserController(App.UserService)
	App.LoadConfiguration()

	routes.RegisterSwaggerRoute(App.EchoInstance)
	routes.RegisterHealthcheckRoute(App.EchoInstance)
	routes.RegisterAuthRoutes(App.EchoInstance, &App.OAuthController)
	routes.RegisterUserRoutes(App.EchoInstance, &App.UserController)

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	App.EchoInstance.Listener = listener

	go func() {
		App.EchoInstance.Start("")
	}()

	return fmt.Sprintf("http://localhost:%d", port)
}

func TestMain(m *testing.M) {
	basePath = InitServer()
	time.Sleep(300 * time.Millisecond)

	code := m.Run()
	os.Exit(code)
}

func makeRequest(httpMethod string, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func TestHealthcheck(t *testing.T) {
	url := fmt.Sprintf("%s/healthcheck/", basePath)
	resp := makeRequest("GET", url, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOAuthToken(t *testing.T) {
	url := fmt.Sprintf("%s/oauth/token", basePath)
	body := strings.NewReader(`{"client_id":"c6cece53","client_secret":"123"}`)
	resp := makeRequest("POST", url, body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOAuthUserInfo(t *testing.T) {
	url := fmt.Sprintf("%s/oauth/userinfo", basePath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer test-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateUser(t *testing.T) {
	url := fmt.Sprintf("%s/api/user/", basePath)
	body := strings.NewReader(`{"username":"test","email":"test@test.com","firstName":"Test","lastName":"User","password":"123"}`)
	resp := makeRequest("POST", url, body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestGetUsers(t *testing.T) {
	url := fmt.Sprintf("%s/api/user/", basePath)
	resp := makeRequest("GET", url, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
