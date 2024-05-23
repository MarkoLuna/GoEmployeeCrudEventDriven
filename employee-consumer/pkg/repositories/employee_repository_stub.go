package repositories

import (
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/constants"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/models"
)

type EmployeeRepositoryStub struct {
}

func NewEmployeeRepositoryStub() EmployeeRepository {
	return &EmployeeRepositoryStub{}
}

func (er EmployeeRepositoryStub) Create(e models.Employee) (*models.Employee, error) {
	return &e, nil
}

func (er EmployeeRepositoryStub) FindAll() ([]models.Employee, error) {

	employeesSlice := make([]models.Employee, 0)

	var employee models.Employee
	employee.Id = "1"
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	employeesSlice = append(employeesSlice, employee)

	var employee2 models.Employee
	employee2.Id = "2"
	employee2.FirstName = "Gerardo"
	employee2.LastName = "Luna"
	employee2.SecondLastName = "Valdez"
	employee2.DateOfBirth = time.Date(1999, time.November, 8, 8, 0, 0, 0, time.UTC)
	employee2.DateOfEmployment = time.Now().UTC()
	employee2.Status = constants.ACTIVE

	employeesSlice = append(employeesSlice, employee2)

	return employeesSlice, nil
}

func (er EmployeeRepositoryStub) FindById(ID string) (models.Employee, error) {

	var employee models.Employee
	employee.Id = ID
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	return employee, nil
}

func (er EmployeeRepositoryStub) DeleteById(ID string) (int64, error) {
	return 1, nil
}

func (er EmployeeRepositoryStub) Update(e models.Employee) (int64, error) {
	return 1, nil
}
