package app

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeService/pkg/controllers"
	"github.com/MarkoLuna/EmployeeService/pkg/dto"
	"github.com/MarkoLuna/EmployeeService/pkg/repositories"
	"github.com/MarkoLuna/EmployeeService/pkg/services"
	"github.com/MarkoLuna/EmployeeService/pkg/services/stubs"
	"github.com/labstack/echo/v4"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MarkoLuna/EmployeeService/pkg/constants"
	"github.com/MarkoLuna/EmployeeService/pkg/models"

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
	App.EmployeeService = services.NewEmployeeService(App.EmployeeRepository)
	App.EmployeeController = controllers.NewEmployeeController(App.EmployeeService)
	App.OAuthService = stubs.NewOAuthServiceStub()

	App.LoadConfiguration()
}

var employeeId = "1"
var e = &dto.EmployeeRequest{
	FirstName:        "Marcos",
	LastName:         "Luna",
	SecondLastName:   "Valdez",
	DateOfBirth:      time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC),
	DateOfEmployment: time.Now().UTC(),
	Status:           constants.ACTIVE,
}

var invalidEmployee = &models.Employee{
	Id:               "1",
	FirstName:        "",
	LastName:         "",
	SecondLastName:   "",
	DateOfBirth:      time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC),
	DateOfEmployment: time.Now().UTC(),
	Status:           "",
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

func TestFindById(t *testing.T) {

	query := `select
				id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status
			from
				employees
			where
				id_employee \\?`

	rows := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"}).
		AddRow(employeeId, e.FirstName, e.LastName, e.SecondLastName, e.DateOfBirth, e.DateOfEmployment, e.Status)

	sqlMock.ExpectQuery(query).WithArgs(employeeId).WillReturnRows(rows)

	url := fmt.Sprintf("%s/api/employee/%s", basePath, employeeId)
	resp := makeRequest("GET", url, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid http status code")

	employeeResponse := models.Employee{}
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, &employeeResponse)

	assert.NotNil(t, employeeResponse)
	assert.NoError(t, err)
}

func TestEmployeeRepositoryImpl_FindByIdError(t *testing.T) {

	query := `select
				id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status
			from
				employees
			where
				id_employee `

	rows := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"})

	sqlMock.ExpectQuery(query).WithArgs(employeeId).WillReturnRows(rows)

	url := fmt.Sprintf("%s/api/employee/%s", basePath, employeeId)
	resp := makeRequest("GET", url, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Invalid http status code")
}

func TestFindAll(t *testing.T) {

	query := `SELECT id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status 
			  FROM employees`

	rows := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"}).
		AddRow(employeeId, e.FirstName, e.LastName, e.SecondLastName, e.DateOfBirth, e.DateOfEmployment, e.Status)

	sqlMock.ExpectQuery(query).WillReturnRows(rows)

	url := fmt.Sprintf("%s/api/employee/", basePath)
	resp := makeRequest("GET", url, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid http status code")

	employeesSlice := make([]models.Employee, 0)
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, &employeesSlice)

	assert.NotEmpty(t, employeesSlice)
	assert.NoError(t, err)
	assert.Len(t, employeesSlice, 1)
}

func TestFindAllError(t *testing.T) {

	query := `SELECT id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status
			  FROM employees`

	rows := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"})

	sqlMock.ExpectQuery(query).WillReturnRows(rows)

	url := fmt.Sprintf("%s/api/employee/", basePath)
	resp := makeRequest("GET", url, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid http status code")

	employeesSlice := make([]models.Employee, 0)
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, &employeesSlice)

	fmt.Println("response Body:", string(body))
	assert.Empty(t, employeesSlice)
	assert.NoError(t, err)
	assert.Len(t, employeesSlice, 0)
}

func TestEmployeeRepositoryImpl_Create(t *testing.T) {

	query := `
	INSERT INTO employees \(
		id_employee\,
		first_name\,
		last_name\,
		second_last_name\,
		date_of_birth\,
		date_of_employment\,
		status\)
	VALUES \(\$1\, \$2\, \$3\, \$4\, \$5\, \$6\, \$7\)
	RETURNING id_employee`

	rows := sqlmock.NewRows([]string{"id_employee"}).AddRow(1)

	sqlMock.ExpectQuery(query).WithArgs(AnyUUID{}, e.FirstName, e.LastName, e.SecondLastName,
		e.DateOfBirth, e.DateOfEmployment, e.Status).WillReturnRows(rows)

	url := fmt.Sprintf("%s/api/employee/", basePath)
	jsonStr, _ := json.Marshal(e)
	resp := makeRequest("POST", url, bytes.NewBuffer(jsonStr))

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "handler returned wrong status code")

	employeeResponse := models.Employee{}
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, &employeeResponse)

	assert.NoError(t, err)
	assert.NotNil(t, employeeResponse)

	assert.Equal(t, e.FirstName, employeeResponse.FirstName, "FirstName employee returned is wrong")
}

func TestEmployeeRepositoryImpl_CreateAndGetInvalidInput(t *testing.T) {

	url := fmt.Sprintf("%s/api/employee/", basePath)
	jsonStr, _ := json.Marshal(invalidEmployee)
	resp := makeRequest("POST", url, bytes.NewBuffer(jsonStr))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "handler returned wrong status code")
}

func TestEmployeeRepositoryImpl_Update(t *testing.T) {

	query_select := `select
				id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status
			from
				employees
			where
				id_employee \\?`

	rows_select := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"}).
		AddRow(employeeId, e.FirstName, e.LastName, e.SecondLastName,
			e.DateOfBirth, e.DateOfEmployment, e.Status)

	sqlMock.ExpectQuery(query_select).WithArgs(employeeId).WillReturnRows(rows_select)

	query_update := `
		UPDATE employees SET
			first_name = \$2\,
			last_name = \$3\,
			second_last_name = \$4\,
			date_of_birth = \$5\,
			date_of_employment = \$6\,
			status = \$7
		WHERE id_employee = \$1;
	`
	sqlMock.ExpectExec(query_update).WithArgs(employeeId, e.FirstName, e.LastName, e.SecondLastName,
		e.DateOfBirth, e.DateOfEmployment, e.Status).WillReturnResult(sqlmock.NewResult(0, 1))

	url := fmt.Sprintf("%s/api/employee/%s", basePath, employeeId)
	jsonStr, _ := json.Marshal(e)
	resp := makeRequest("PUT", url, bytes.NewBuffer(jsonStr))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "handler returned wrong status code")

	employeeResponse := models.Employee{}
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, &employeeResponse)

	assert.NoError(t, err)
	assert.NotNil(t, employeeResponse)

	assert.Equal(t, employeeId, employeeResponse.Id, "id employee returned is wrong")
}

func TestEmployeeRepositoryImpl_UpdateErr(t *testing.T) {

	query_select := `select
				id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status
			from
				employees
			where
				id_employee \\?`

	rows_select := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"})

	sqlMock.ExpectQuery(query_select).WithArgs(employeeId).WillReturnRows(rows_select)

	url := fmt.Sprintf("%s/api/employee/%s", basePath, employeeId)
	jsonStr, _ := json.Marshal(e)
	resp := makeRequest("PUT", url, bytes.NewBuffer(jsonStr))

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "handler returned wrong status code")
}

func TestEmployeeRepositoryImpl_DeleteById(t *testing.T) {

	query := `DELETE FROM employees WHERE id_employee = \$1\;`
	sqlMock.ExpectExec(query).WithArgs(employeeId).WillReturnResult(sqlmock.NewResult(0, 1))

	url := fmt.Sprintf("%s/api/employee/%s", basePath, employeeId)
	resp := makeRequest("DELETE", url, nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "handler returned wrong status code")
}

func TestEmployeeRepositoryStub_DeleteByIdError(t *testing.T) {

	query := `DELETE FROM employees WHERE id_employee = \$1\;`
	sqlMock.ExpectExec(query).WithArgs(employeeId).WillReturnResult(sqlmock.NewResult(0, 0))

	url := fmt.Sprintf("%s/api/employee/%s", basePath, employeeId)
	resp := makeRequest("DELETE", url, nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "handler returned wrong status code")
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
