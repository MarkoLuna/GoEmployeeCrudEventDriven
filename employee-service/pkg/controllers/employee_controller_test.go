package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/labstack/echo/v4"

	"github.com/MarkoLuna/EmployeeService/pkg/clients"
	"github.com/MarkoLuna/EmployeeService/pkg/constants"
	"github.com/MarkoLuna/EmployeeService/pkg/models"
	"github.com/MarkoLuna/EmployeeService/pkg/services"
	"github.com/MarkoLuna/EmployeeService/pkg/services/stubs"
	"github.com/stretchr/testify/assert"
)

var invalidEmployee = &models.Employee{
	Id:               "1",
	FirstName:        "",
	LastName:         "",
	SecondLastName:   "",
	DateOfBirth:      time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC),
	DateOfEmployment: time.Now().UTC(),
	Status:           "",
}

var validEmployee = &models.Employee{
	Id:               "1",
	FirstName:        "Marcos",
	LastName:         "Luna",
	SecondLastName:   "Valdez",
	DateOfBirth:      time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC),
	DateOfEmployment: time.Now().UTC(),
	Status:           constants.ACTIVE,
}

func TestEmployeeController_GetEmployeesEmployees(t *testing.T) {

	employees := make([]models.Employee, 0)
	employees = append(employees, createEmployee1())
	employees = append(employees, createEmployee2())

	employeeClient := clients.NewEmployeeConsumerServiceStubFromData(employees, nil)
	kafkaProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	req, err := http.NewRequest(http.MethodGet, "/api/employee/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	c := e.NewContext(req, rr)
	if assert.NoError(t, employeeController.GetEmployees(c)) {

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		assert.NotEqual(t, 0, len(rr.Body.String()), "handler returned unexpected body: got empty")

		employeesSlice := make([]models.Employee, 0)
		err = json.Unmarshal(rr.Body.Bytes(), &employeesSlice)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(employeesSlice), "handler returned unexpected body: got empty")

		employee1 := employeesSlice[0]
		employee2 := employeesSlice[1]

		assert.NotNil(t, employee1)
		assert.NotNil(t, employee2)

		fmt.Println(employee1)
		fmt.Println(employee2)

		assert.Equal(t, "1", employee1.Id, "id employee returned is wrong")
		assert.Equal(t, "2", employee2.Id, "id employee returned is wrong")
	}
}

func TestEmployeeController_CreateEmployeeEmployee(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	kafkaProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	jsonStr, _ := json.Marshal(validEmployee)
	req, err := http.NewRequest(http.MethodPost, "/api/employee/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)

	if assert.NoError(t, employeeController.CreateEmployee(c)) {
		assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
	}
}

func TestEmployeeController_CreateEmployeeEmployeeThenBadRequest(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	kafkaProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	jsonStr, _ := json.Marshal(invalidEmployee)
	req, err := http.NewRequest(http.MethodPost, "/api/employee/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)

	if assert.NoError(t, employeeController.CreateEmployee(c)) {
		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	}
}

func TestEmployeeController_GetEmployeeByIdEmployee(t *testing.T) {

	employees := make([]models.Employee, 0)
	employees = append(employees, *validEmployee)

	employeeClient := clients.NewEmployeeConsumerServiceStubFromData(employees, nil)
	kafkaProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	req, err := http.NewRequest(http.MethodGet, "/api/employee/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)

	c.SetPath("/api/employee/:employeeId")
	c.SetParamNames("employeeId")
	c.SetParamValues("1")

	if assert.NoError(t, employeeController.GetEmployeeById(c)) {
		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		assert.NotEqual(t, 0, len(rr.Body.String()), "handler returned unexpected body: got empty")

		employeeResponse := models.Employee{}
		err = json.Unmarshal(rr.Body.Bytes(), &employeeResponse)

		assert.NoError(t, err)

		assert.NotNil(t, employeeResponse)
		fmt.Println(employeeResponse)

		assert.Equal(t, "1", employeeResponse.Id, "id employee returned is wrong")
	}
}

func TestEmployeeController_GetEmployeeByIdEmployeeThenNotFound(t *testing.T) {

	employees := make([]models.Employee, 0)
	employees = append(employees, *validEmployee)

	employeeClient := clients.NewEmployeeConsumerServiceStubFromError(errors.New("employee not Found"))
	kafkaProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	req, err := http.NewRequest(http.MethodGet, "/api/employee/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)

	c.SetPath("/api/employee/:employeeId")
	c.SetParamNames("employeeId")
	c.SetParamValues("1")

	if assert.NoError(t, employeeController.GetEmployeeById(c)) {
		assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
	}
}

func TestEmployeeController_UpdateEmployee(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	kafkaProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	jsonStr, _ := json.Marshal(validEmployee)
	req, err := http.NewRequest(http.MethodPut, "/api/employee/1", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)

	c.SetPath("/api/employee/:employeeId")
	c.SetParamNames("employeeId")
	c.SetParamValues("1")

	if assert.NoError(t, employeeController.UpdateEmployee(c)) {
		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	}
}

func TestEmployeeController_UpdateEmployeeThenNotFound(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStubFromError(errors.New("employee not Found"))
	kafkaProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	jsonStr, _ := json.Marshal(validEmployee)
	req, err := http.NewRequest(http.MethodPut, "/api/employee/1", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)

	c.SetPath("/api/employee/:employeeId")
	c.SetParamNames("employeeId")
	c.SetParamValues("1")

	if assert.NoError(t, employeeController.UpdateEmployee(c)) {
		assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
	}
}

func TestEmployeeController_DeleteEmployee(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	kafkaProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	req, err := http.NewRequest(http.MethodDelete, "/api/employee/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)

	c.SetPath("/api/employee/:employeeId")
	c.SetParamNames("employeeId")
	c.SetParamValues("1")

	if assert.NoError(t, employeeController.DeleteEmployee(c)) {
		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	}

}

func TestDeleteByIdError(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	kafkaProducerService := stubs.NewKafkaProducerServiceStubFromError(errors.New("failed to delete"))
	employeeService := services.NewEmployeeService(employeeClient, kafkaProducerService)
	employeeController := NewEmployeeController(employeeService)

	e := echo.New()

	req, err := http.NewRequest(http.MethodDelete, "/api/employee/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)

	c.SetPath("/api/employee/:employeeId")
	c.SetParamNames("employeeId")
	c.SetParamValues("1")

	if assert.NoError(t, employeeController.DeleteEmployee(c)) {
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")
	}
}

func createEmployee1() models.Employee {
	var employee models.Employee
	employee.Id = "1"
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	return employee
}

func createEmployee2() models.Employee {

	var employee2 models.Employee
	employee2.Id = "2"
	employee2.FirstName = "Gerardo"
	employee2.LastName = "Luna"
	employee2.SecondLastName = "Valdez"
	employee2.DateOfBirth = time.Date(1999, time.November, 8, 8, 0, 0, 0, time.UTC)
	employee2.DateOfEmployment = time.Now().UTC()
	employee2.Status = constants.ACTIVE

	return employee2
}
