package repositories

import (
	"github.com/MarkoLuna/EmployeeConsumer/pkg/models"
)

type EmployeeRepository interface {
	Create(e models.Employee) (*models.Employee, error)

	FindAll() ([]models.Employee, error)

	FindById(ID string) (models.Employee, error)

	DeleteById(ID string) (int64, error)

	Update(e models.Employee) (int64, error)
}
