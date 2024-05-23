package utils

import (
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/constants"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/models"
)

func TestValidEmployee(t *testing.T) {

	var employee dto.EmployeeRequest
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Now().UTC()
	employee.Status = constants.ACTIVE

	v := CreateValidator()
	err1 := v.Struct(employee)
	if err1 != nil {
		t.Error("the employee should be valid")
	}
}

func TestParseEmployeeRequestWithInvalidData(t *testing.T) {

	employee := &models.Employee{}

	v := CreateValidator()
	err1 := v.Struct(employee)
	if err1 == nil {
		t.Error("the body should not be valid")
	}
}
