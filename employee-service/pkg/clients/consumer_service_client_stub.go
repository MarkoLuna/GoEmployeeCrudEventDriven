package clients

import (
	"github.com/MarkoLuna/EmployeeService/pkg/models"
)

type EmployeeConsumerServiceClientStub struct {
	employees []models.Employee
	err       error
}

func NewEmployeeConsumerServiceStub() EmployeeConsumerServiceClient {
	return &EmployeeConsumerServiceClientStub{employees: nil, err: nil}
}

func NewEmployeeConsumerServiceStubFromData(employees []models.Employee, err error) EmployeeConsumerServiceClient {
	return EmployeeConsumerServiceClientStub{employees: employees, err: err}
}

func NewEmployeeConsumerServiceStubFromEmployee(e models.Employee) EmployeeConsumerServiceClient {
	employees := make([]models.Employee, 1)
	employees = append(employees, e)
	service := &EmployeeConsumerServiceClientStub{employees: employees, err: nil}

	return service
}

func NewEmployeeConsumerServiceStubFromEmployes(employees []models.Employee) EmployeeConsumerServiceClient {
	return &EmployeeConsumerServiceClientStub{employees: employees, err: nil}
}

func NewEmployeeConsumerServiceStubFromError(err error) EmployeeConsumerServiceClient {
	return &EmployeeConsumerServiceClientStub{employees: nil, err: err}
}

func (es EmployeeConsumerServiceClientStub) Create(e models.Employee) (models.Employee, error) {
	if len(es.employees) == 0 {
		return models.Employee{}, es.err
	}

	return es.employees[0], es.err
}

func (es EmployeeConsumerServiceClientStub) FindAll() ([]models.Employee, error) {
	return es.employees, es.err
}

func (es EmployeeConsumerServiceClientStub) FindById(ID string) (models.Employee, error) {
	if len(es.employees) == 0 {
		return models.Employee{}, es.err
	}

	return es.employees[0], es.err
}

func (es EmployeeConsumerServiceClientStub) DeleteById(ID string) error {
	return es.err
}

func (es EmployeeConsumerServiceClientStub) Update(e models.Employee) (models.Employee, error) {
	if len(es.employees) == 0 {
		return models.Employee{}, es.err
	}

	return es.employees[0], es.err
}
