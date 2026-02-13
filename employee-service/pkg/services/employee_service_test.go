package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeService/pkg/clients"
	"github.com/MarkoLuna/EmployeeService/pkg/constants"
	"github.com/MarkoLuna/EmployeeService/pkg/dto"
	"github.com/MarkoLuna/EmployeeService/pkg/models"
	"github.com/MarkoLuna/EmployeeService/pkg/services/stubs"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeService_GetEmployeesEmployees(t *testing.T) {

	employees := make([]models.Employee, 0)
	employees = append(employees, createEmployee1())
	employees = append(employees, createEmployee2())

	employeeClient := clients.NewEmployeeConsumerServiceStubFromData(employees, nil)
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(employeeClient, employeeProducerService)

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

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(employeeClient, employeeProducerService)

	var employee dto.EmployeeRequest
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	_, err := employeeService.CreateEmployee(employee)

	assert.NoError(t, err)
}

func TestEmployeeService_GetEmployeeByIdEmployee(t *testing.T) {

	employees := make([]models.Employee, 0)
	employees = append(employees, createEmployee1())

	employeeClient := clients.NewEmployeeConsumerServiceStubFromData(employees, nil)
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(employeeClient, employeeProducerService)

	employeeResponse, err := employeeService.GetEmployeeById("1")

	assert.NoError(t, err)
	assert.NotNil(t, employeeResponse)
	fmt.Println(employeeResponse)

	assert.Equal(t, "1", employeeResponse.Id, "id employee returned is wrong")
}

func TestEmployeeService_UpdateEmployee(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(employeeClient, employeeProducerService)

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
}

func TestEmployeeService_DeleteEmployee(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(employeeClient, employeeProducerService)

	err := employeeService.DeleteEmployeeById("1")
	assert.NoError(t, err)
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
