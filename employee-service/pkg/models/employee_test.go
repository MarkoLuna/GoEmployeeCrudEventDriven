package models

import (
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeService/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestEmployee_ToString(t *testing.T) {
	var employee Employee
	employee.Id = "1"
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.Status = constants.ACTIVE

	jsonEmployee := `{"id":"1","firstName":"Marcos","lastName":"Luna","secondLastName":"Valdez","dateOfBirth":"1994-04-25T08:00:00Z","dateOfEmployment":"1994-04-25T08:00:00Z","status":"ACTIVE"}`

	assert.Equal(t, jsonEmployee, employee.ToString())
}
