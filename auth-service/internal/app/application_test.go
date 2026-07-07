package app

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/MarkoLuna/AuthService/internal/controllers"
	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/MarkoLuna/AuthService/internal/services/impl"
	"github.com/MarkoLuna/AuthService/internal/services/stubs"
	"github.com/labstack/echo/v4"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	App          = Application{}
	basePath     = "http://localhost:8082"
	dbConnection *sql.DB
	sqlMock      sqlmock.Sqlmock
)

func InitServer(dbConnection *sql.DB) {
	App.EchoInstance = echo.New()

	App.UserRepository = repositories.NewUserRepository(dbConnection, false)
	App.UserService = impl.NewUserService(App.UserRepository)
	App.AuthStrategy = stubs.NewAuthStrategyStub()

	App.OAuthController = controllers.NewOAuthController(App.AuthStrategy)
	App.UserController = controllers.NewUserController(App.UserService)

	App.LoadConfiguration()
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestMain(m *testing.M) {
	db, mock := NewMock()
	sqlMock = mock
	dbConnection = db
	InitServer(dbConnection)
	os.Setenv("SERVER_SSL_ENABLED", "false")
	go App.Run()

	code := m.Run()
	shutdown()
	os.Exit(code)
}

func shutdown() {
	defer dbConnection.Close()
}

func TestHealthCheck(t *testing.T) {
	url := fmt.Sprintf("%s/healthcheck/", basePath)
	resp := makeRequest("GET", url, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "handler returned wrong status code")
}

func makeRequest(httpMethod string, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(httpMethod, url, body)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}
