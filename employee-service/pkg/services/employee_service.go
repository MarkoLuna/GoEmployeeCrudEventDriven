package services

import (
	"errors"

	"github.com/MarkoLuna/EmployeeService/pkg/clients"
	"github.com/MarkoLuna/EmployeeService/pkg/dto"
	"github.com/MarkoLuna/EmployeeService/pkg/models"
)

type EmployeeService struct {
	employeeClient          clients.EmployeeConsumerServiceClient
	employeeProducerService KafkaProducerService
}

func NewEmployeeService(employeeClient clients.EmployeeConsumerServiceClient,
	employeeProducerService KafkaProducerService) EmployeeService {
	return EmployeeService{employeeClient: employeeClient,
		employeeProducerService: employeeProducerService}
}

func (eSrv EmployeeService) CreateEmployee(employeeRequest dto.EmployeeRequest) (*models.Employee, error) {

	err := eSrv.employeeProducerService.SendUpsert(employeeRequest)
	return nil, err
}

func (eSrv EmployeeService) GetEmployees() ([]models.Employee, error) {
	employees, err := eSrv.employeeClient.FindAll()
	return employees, err
}

func (eSrv EmployeeService) GetEmployeeById(employeeId string) (models.Employee, error) {
	employeeDetails, err := eSrv.employeeClient.FindById(employeeId)
	return employeeDetails, err
}

func (eSrv EmployeeService) UpdateEmployee(employeeId string, employee dto.EmployeeRequest) (models.Employee, error) {

	_, err := eSrv.employeeClient.FindById(employeeId)
	if err == nil {

		err := eSrv.employeeProducerService.SendUpsert(employee)
		return models.Employee{}, err
	} else {
		return models.Employee{}, errors.New("employee not Found")
	}
}

func (eSrv EmployeeService) DeleteEmployeeById(employeeId string) error {
	err := eSrv.employeeProducerService.SendDelete(employeeId)
	return err
}
