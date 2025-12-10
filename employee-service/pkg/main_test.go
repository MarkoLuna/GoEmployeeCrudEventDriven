package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MarkoLuna/EmployeeService/pkg/constants"
	"github.com/MarkoLuna/EmployeeService/pkg/models"
	"github.com/MarkoLuna/EmployeeService/pkg/repositories"
	"github.com/MarkoLuna/EmployeeService/pkg/services/stubs"

	"github.com/stretchr/testify/assert"
)

var (
	basePath = "http://localhost:8080"
	sqlMock  sqlmock.Sqlmock
)

func InitServer(dbConnection *sql.DB) {
	App.DbConnection = dbConnection
	App.EmployeeRepository = repositories.NewEmployeeRepository(App.DbConnection, false)
	App.OAuthService = stubs.NewOAuthServiceStub()
	go main()
}

var e = &models.Employee{
	Id:               "1",
	FirstName:        "Marcos",
	LastName:         "Luna",
	SecondLastName:   "Valdez",
	DateOfBirth:      time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC),
	DateOfEmployment: time.Now().UTC(),
	Status:           constants.ACTIVE,
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
	InitServer(db)

	go func() {
		code := m.Run()
		os.Exit(code)
	}()
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
		AddRow(e.Id, e.FirstName, e.LastName, e.SecondLastName, e.DateOfBirth, e.DateOfEmployment, e.Status)

	sqlMock.ExpectQuery(query).WithArgs(e.Id).WillReturnRows(rows)

	url := fmt.Sprintf("%s/api/employee/%s", basePath, e.Id)
	resp := makeRequest("GET", url, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid http status code")

	employeeResponse := models.Employee{}
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, &employeeResponse)

	assert.NotNil(t, employeeResponse)
	assert.NoError(t, err)
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
		AddRow(e.Id, e.FirstName, e.LastName, e.SecondLastName, e.DateOfBirth, e.DateOfEmployment, e.Status)

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
