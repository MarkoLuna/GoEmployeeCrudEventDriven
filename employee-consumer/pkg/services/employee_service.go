package services

import (
	"errors"
	"log"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/models"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/repositories"
	"github.com/google/uuid"
)

type EmployeeService struct {
	employeeRepository repositories.EmployeeRepository
}

func NewEmployeeService(employeeRepository repositories.EmployeeRepository) EmployeeService {
	return EmployeeService{employeeRepository}
}

func (eSrv EmployeeService) CreateEmployee(employeeRequest dto.EmployeeRequest) (*models.Employee, error) {

	var employee models.Employee
	employee.Id = uuid.New().String()
	employee.FirstName = employeeRequest.FirstName
	employee.LastName = employeeRequest.LastName
	employee.SecondLastName = employeeRequest.SecondLastName
	employee.DateOfBirth = employeeRequest.DateOfBirth
	employee.DateOfEmployment = employeeRequest.DateOfEmployment
	employee.Status = employeeRequest.Status

	log.Println("employee: " + employee.ToString())
	e, err := eSrv.employeeRepository.Create(employee)

	return e, err
}

func (eSrv EmployeeService) GetEmployees() ([]models.Employee, error) {
	employees, err := eSrv.employeeRepository.FindAll()
	return employees, err
}

func (eSrv EmployeeService) GetEmployeeById(employeeId string) (models.Employee, error) {
	employeeDetails, err := eSrv.employeeRepository.FindById(employeeId)
	return employeeDetails, err
}

func (eSrv EmployeeService) UpdateEmployee(employeeId string, employee dto.EmployeeRequest) (models.Employee, error) {

	employeeDetails, err := eSrv.employeeRepository.FindById(employeeId)
	if err == nil {
		employeeDetails.FirstName = employee.FirstName
		employeeDetails.LastName = employee.LastName
		employeeDetails.SecondLastName = employee.SecondLastName
		employeeDetails.DateOfBirth = employee.DateOfBirth
		employeeDetails.DateOfEmployment = employee.DateOfEmployment
		employeeDetails.Status = employee.Status

		log.Println("employee: " + employeeDetails.ToString())

		count, _ := eSrv.employeeRepository.Update(employeeDetails)
		if count > 0 {
			return employeeDetails, nil
		} else {
			return models.Employee{}, errors.New("employee not Found")
		}
	} else {
		return models.Employee{}, errors.New("employee not Found")
	}
}

func (eSrv EmployeeService) DeleteEmployeeById(employeeId string) error {
	count, _ := eSrv.employeeRepository.DeleteById(employeeId)
	if count > 0 {
		return nil
	} else {
		return errors.New("employee not Found")
	}
}
