package app

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/MarkoLuna/EmployeeService/pkg/clients"
	"github.com/MarkoLuna/EmployeeService/pkg/controllers"
	"github.com/MarkoLuna/EmployeeService/pkg/repositories"
	"github.com/MarkoLuna/EmployeeService/pkg/services"
	"github.com/MarkoLuna/EmployeeService/pkg/services/stubs"
	"github.com/labstack/echo/v4"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

var (
	App      = Application{}
	basePath = "http://localhost:8080"

	dbConnection *sql.DB
	sqlMock      sqlmock.Sqlmock
)

func InitServer(db_connection *sql.DB) {
	App.EchoInstance = echo.New()
	App.DbConnection = db_connection
	App.EmployeeRepository = repositories.NewEmployeeRepository(App.DbConnection, false)
	App.EmployeeConsumerServiceClient = clients.NewEmployeeConsumerServiceStub()
	App.KafkaProducerService = stubs.NewKafkaProducerServiceStub()
	App.EmployeeService = services.NewEmployeeService(App.EmployeeConsumerServiceClient, App.KafkaProducerService)
	App.EmployeeController = controllers.NewEmployeeController(App.EmployeeService)
	App.OAuthService = stubs.NewOAuthServiceStub()

	App.LoadConfiguration()
}

type AnyUUID struct{}

// Match satisfies sqlmock.Argument interface
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
