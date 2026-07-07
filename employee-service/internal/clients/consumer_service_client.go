package clients

import (
	"context"

	"github.com/MarkoLuna/EmployeeService/internal/models"
)

type EmployeeConsumerServiceClient interface {
	Create(ctx context.Context, e models.Employee) (models.Employee, error)

	FindAll(ctx context.Context) ([]models.Employee, error)

	FindById(ctx context.Context, ID string) (models.Employee, error)

	DeleteById(ctx context.Context, ID string) error

	Update(ctx context.Context, e models.Employee) (models.Employee, error)
}
