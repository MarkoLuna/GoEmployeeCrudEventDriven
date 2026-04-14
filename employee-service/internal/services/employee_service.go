package services

import (
	"errors"

	"github.com/MarkoLuna/EmployeeService/internal/clients"
	"github.com/MarkoLuna/EmployeeService/internal/models"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
)

type EmployeeService struct {
	clientBuilder           *clients.EmployeeConsumerServiceClientBuilder
	employeeProducerService KafkaProducerService
}

func NewEmployeeService(clientBuilder *clients.EmployeeConsumerServiceClientBuilder,
	employeeProducerService KafkaProducerService) EmployeeService {
	return EmployeeService{
		clientBuilder:           clientBuilder,
		employeeProducerService: employeeProducerService,
	}
}

func (eSrv EmployeeService) getClient(jwt string) clients.EmployeeConsumerServiceClient {
	return eSrv.clientBuilder.WithJwtToken(jwt).Build()
}

func (eSrv EmployeeService) CreateEmployee(jwt string, employeeRequest dto.EmployeeRequest) (*models.Employee, error) {

	employeeMessage := dto.EmployeeMessage{EmployeeInfo: employeeRequest}
	err := eSrv.employeeProducerService.SendInsert(employeeMessage)
	return nil, err
}

func (eSrv EmployeeService) GetEmployees(jwt string) ([]models.Employee, error) {
	client := eSrv.getClient(jwt)
	employees, err := client.FindAll()
	return employees, err
}

func (eSrv EmployeeService) GetEmployeeById(jwt string, employeeId string) (models.Employee, error) {
	client := eSrv.getClient(jwt)
	employeeDetails, err := client.FindById(employeeId)
	return employeeDetails, err
}

func (eSrv EmployeeService) UpdateEmployee(jwt string, employeeId string, employee dto.EmployeeRequest) (models.Employee, error) {
	client := eSrv.getClient(jwt)
	currentEmployee, err := client.FindById(employeeId)
	if err == nil {

		employeeMessage := dto.EmployeeMessage{ID: currentEmployee.Id, EmployeeInfo: employee}
		err := eSrv.employeeProducerService.SendUpdate(employeeMessage)
		return models.Employee{}, err
	} else {
		return models.Employee{}, errors.New("employee not Found")
	}
}

func (eSrv EmployeeService) DeleteEmployeeById(jwt string, employeeId string) error {
	err := eSrv.employeeProducerService.SendDelete(employeeId)
	return err
}
