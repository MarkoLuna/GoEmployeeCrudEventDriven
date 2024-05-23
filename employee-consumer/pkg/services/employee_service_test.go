package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/constants"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/repositories"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeService_GetEmployeesEmployees(t *testing.T) {

	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := NewEmployeeService(employeeRepository)

	employeesSlice, err := employeeService.GetEmployees()

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

func TestEmployeeService_CreateEmployeeEmployee(t *testing.T) {

	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := NewEmployeeService(employeeRepository)

	var employee dto.EmployeeRequest
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	employeeResponse, err := employeeService.CreateEmployee(employee)

	assert.NoError(t, err)
	assert.NotNil(t, employeeResponse)
	fmt.Println(employeeResponse)

	assert.Equal(t, employee.FirstName, employeeResponse.FirstName, "FirstName employee returned is wrong")
}

func TestEmployeeService_GetEmployeeByIdEmployee(t *testing.T) {

	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := NewEmployeeService(employeeRepository)

	employeeResponse, err := employeeService.GetEmployeeById("1")

	assert.NoError(t, err)
	assert.NotNil(t, employeeResponse)
	fmt.Println(employeeResponse)

	assert.Equal(t, "1", employeeResponse.Id, "id employee returned is wrong")
}

func TestEmployeeService_UpdateEmployee(t *testing.T) {

	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := NewEmployeeService(employeeRepository)

	var employee dto.EmployeeRequest
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	employeeResponse, err := employeeService.UpdateEmployee("1", employee)

	assert.NoError(t, err)
	assert.NotNil(t, employeeResponse)
	fmt.Println(employeeResponse)

	assert.Equal(t, "1", employeeResponse.Id, "id employee returned is wrong")
}

func TestEmployeeService_DeleteEmployee(t *testing.T) {

	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := NewEmployeeService(employeeRepository)

	err := employeeService.DeleteEmployeeById("1")
	assert.NoError(t, err)
}
