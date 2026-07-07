package app

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/MarkoLuna/EmployeeConsumer/internal/controllers"
	"github.com/MarkoLuna/EmployeeConsumer/internal/repositories"
	"github.com/MarkoLuna/EmployeeConsumer/internal/services/impl"
	"github.com/labstack/echo/v4"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

var (
	App      = Application{}
	basePath = "http://localhost:8081"

	dbConnection *sql.DB
	sqlMock      sqlmock.Sqlmock
)

func InitServer(dbConnection *sql.DB) {
	App.EchoInstance = echo.New()
	App.DbConnection = dbConnection
	App.EmployeeRepository = repositories.NewEmployeeRepository(App.DbConnection, false)
	App.EmployeeService = impl.NewEmployeeService(App.EmployeeRepository)
	App.EmployeeController = controllers.NewEmployeeController(App.EmployeeService)

	App.LoadConfiguration()
}

type AnyUUID struct{}

func (a AnyUUID) Match(v driver.Value) bool {
	value, ok := v.(string)
	if ok {
		r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
		fmt.Println("value: " + value)
		return r.MatchString(value)
	}
	return false
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
	os.Setenv("OAUTH_ENABLED", "false")
	ctx, cancel := context.WithCancel(context.Background())
	go App.Run(ctx)

	code := m.Run()
	cancel()
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

func TestHealthCheckWithSsl(t *testing.T) {
	t.Skip("Skipping test: need to refactor this unit tests logic")
	os.Setenv("SERVER_PORT", "8082")
	os.Setenv("SERVER_SSL_ENABLED", "true")
	defer os.Unsetenv("SERVER_PORT")
	defer os.Unsetenv("SERVER_SSL_ENABLED")

	host := "https://localhost:8082"
	url := fmt.Sprintf("%s/healthcheck/", host)
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
