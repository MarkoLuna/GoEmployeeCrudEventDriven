package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeService/internal/clients"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/constants"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/MarkoLuna/EmployeeService/internal/models"
	"github.com/MarkoLuna/EmployeeService/internal/services/stubs"
	"github.com/stretchr/testify/assert"
)

const testJwt = "test-token"

func TestEmployeeService_GetEmployeesEmployees(t *testing.T) {

	employees := make([]models.Employee, 0)
	employees = append(employees, createEmployee1())
	employees = append(employees, createEmployee2())

	employeeClient := clients.NewEmployeeConsumerServiceStubFromData(employees, nil)
	clientBuilder := clients.NewEmployeeConsumerServiceClientBuilder().WithCustomInstance(employeeClient)
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(clientBuilder, employeeProducerService)

	employeesSlice, err := employeeService.GetEmployees(testJwt)

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
	clientBuilder := clients.NewEmployeeConsumerServiceClientBuilder().WithCustomInstance(employeeClient)
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(clientBuilder, employeeProducerService)

	var employee dto.EmployeeRequest
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	_, err := employeeService.CreateEmployee(testJwt, employee)

	assert.NoError(t, err)
}

func TestEmployeeService_GetEmployeeByIdEmployee(t *testing.T) {

	employees := make([]models.Employee, 0)
	employees = append(employees, createEmployee1())

	employeeClient := clients.NewEmployeeConsumerServiceStubFromData(employees, nil)
	clientBuilder := clients.NewEmployeeConsumerServiceClientBuilder().WithCustomInstance(employeeClient)
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(clientBuilder, employeeProducerService)

	employeeResponse, err := employeeService.GetEmployeeById(testJwt, "1")

	assert.NoError(t, err)
	assert.NotNil(t, employeeResponse)
	fmt.Println(employeeResponse)

	assert.Equal(t, "1", employeeResponse.Id, "id employee returned is wrong")
}

func TestEmployeeService_UpdateEmployee(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	clientBuilder := clients.NewEmployeeConsumerServiceClientBuilder().WithCustomInstance(employeeClient)
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(clientBuilder, employeeProducerService)

	var employee dto.EmployeeRequest
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	employeeResponse, err := employeeService.UpdateEmployee(testJwt, "1", employee)

	assert.NoError(t, err)
	assert.NotNil(t, employeeResponse)
}

func TestEmployeeService_DeleteEmployee(t *testing.T) {

	employeeClient := clients.NewEmployeeConsumerServiceStub()
	clientBuilder := clients.NewEmployeeConsumerServiceClientBuilder().WithCustomInstance(employeeClient)
	employeeProducerService := stubs.NewKafkaProducerServiceStub()
	employeeService := NewEmployeeService(clientBuilder, employeeProducerService)

	err := employeeService.DeleteEmployeeById(testJwt, "1")
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
