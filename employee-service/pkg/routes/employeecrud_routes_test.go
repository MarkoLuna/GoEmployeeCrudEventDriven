package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeService/pkg/clients"
	"github.com/MarkoLuna/EmployeeService/pkg/constants"
	"github.com/MarkoLuna/EmployeeService/pkg/controllers"
	"github.com/MarkoLuna/EmployeeService/pkg/dto"
	"github.com/MarkoLuna/EmployeeService/pkg/services"
	"github.com/MarkoLuna/EmployeeService/pkg/services/stubs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterEmployeeStoreRoutes(t *testing.T) {
	echoInstance := echo.New()

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, employeeProducerService)
	employeeController := controllers.NewEmployeeController(employeeService)

	RegisterEmployeeStoreRoutes(echoInstance, &employeeController)

	var employee dto.EmployeeRequest
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	jsonStr, _ := json.Marshal(employee)

	tables := []struct {
		method  string
		path    string
		body    io.Reader
		status  int
		handler func(c echo.Context) error
	}{
		{"GET", "/api/employee/", nil, http.StatusOK, employeeController.GetEmployees},
		{"POST", "/api/employee/", bytes.NewBuffer(jsonStr), http.StatusCreated, employeeController.CreateEmployee},
		{"GET", "/api/employee/1", nil, http.StatusOK, employeeController.GetEmployeeById},
		{"PUT", "/api/employee/1", bytes.NewBuffer(jsonStr), http.StatusOK, employeeController.UpdateEmployee},
		{"DELETE", "/api/employee/1", nil, http.StatusOK, employeeController.DeleteEmployee},
	}

	for _, table := range tables {
		req, err := http.NewRequest(table.method, table.path, table.body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		c := echoInstance.NewContext(req, rr)

		if assert.NoError(t, table.handler(c)) {
			assert.Equal(t, table.status, rr.Code, "handler returned wrong status code")
		}
	}
}

func TestRegisterHealthcheckRoute(t *testing.T) {
	echoInstance := echo.New()
	RegisterHealthcheckRoute(echoInstance)

	tables := []struct {
		method   string
		path     string
		response string
		status   int
	}{
		{"GET", "/healthcheck/", `OK`, http.StatusOK},
	}

	for _, table := range tables {
		req, err := http.NewRequest(table.method, table.path, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		c := echoInstance.NewContext(req, rr)

		if assert.NoError(t, controllers.HealthCheckHandler(c)) {
			assert.Equal(t, table.response, rr.Body.String(), "handler returned unexpected body")
		}
	}
}
