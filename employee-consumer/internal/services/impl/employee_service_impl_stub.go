package impl

import (
	"time"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/constants"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/MarkoLuna/EmployeeConsumer/internal/models"
	"github.com/google/uuid"
)

type EmployeeServiceStub struct {
}

func NewEmployeeServiceStub() EmployeeServiceStub {
	return EmployeeServiceStub{}
}

func (eSrv EmployeeServiceStub) CreateEmployee(employeeRequest dto.EmployeeRequest) (*models.Employee, error) {

	var employee models.Employee
	employee.Id = uuid.New().String()
	employee.FirstName = employeeRequest.FirstName
	employee.LastName = employeeRequest.LastName
	employee.SecondLastName = employeeRequest.SecondLastName
	employee.DateOfBirth = employeeRequest.DateOfBirth
	employee.DateOfEmployment = employeeRequest.DateOfEmployment
	employee.Status = employeeRequest.Status

	return &employee, nil
}

func (eSrv EmployeeServiceStub) GetEmployees() ([]models.Employee, error) {

	e, err := eSrv.GetEmployeeById("id")
	employeesSlice := make([]models.Employee, 0)
	employeesSlice = append(employeesSlice, e)

	return employeesSlice, err
}

func (eSrv EmployeeServiceStub) GetEmployeeById(employeeId string) (models.Employee, error) {
	employeeDetails := models.Employee{}
	employeeDetails.FirstName = "first name"
	employeeDetails.LastName = "last name"
	employeeDetails.SecondLastName = "second last name"
	employeeDetails.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employeeDetails.DateOfEmployment = time.Now().UTC()
	employeeDetails.Status = constants.ACTIVE

	return employeeDetails, nil
}

func (eSrv EmployeeServiceStub) UpdateEmployee(employeeId string, employee dto.EmployeeRequest) (models.Employee, error) {

	employeeDetails := models.Employee{}
	employeeDetails.FirstName = employee.FirstName
	employeeDetails.LastName = employee.LastName
	employeeDetails.SecondLastName = employee.SecondLastName
	employeeDetails.DateOfBirth = employee.DateOfBirth
	employeeDetails.DateOfEmployment = employee.DateOfEmployment
	employeeDetails.Status = employee.Status

	return employeeDetails, nil
}

func (eSrv EmployeeServiceStub) DeleteEmployeeById(employeeId string) error {
	return nil
}
