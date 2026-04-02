package services

import (
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/MarkoLuna/EmployeeConsumer/internal/models"
)

type EmployeeService interface {
	CreateEmployee(employeeRequest dto.EmployeeRequest) (*models.Employee, error)
	GetEmployees() ([]models.Employee, error)
	GetEmployeeById(employeeId string) (models.Employee, error)
	UpdateEmployee(employeeId string, employee dto.EmployeeRequest) (*models.Employee, error)
	DeleteEmployeeById(employeeId string) error
}
