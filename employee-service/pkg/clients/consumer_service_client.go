package clients

import "github.com/MarkoLuna/EmployeeService/pkg/models"

type EmployeeConsumerServiceClient interface {
	Create(e models.Employee) (models.Employee, error)

	FindAll() ([]models.Employee, error)

	FindById(ID string) (models.Employee, error)

	DeleteById(ID string) error

	Update(e models.Employee) (models.Employee, error)
}
