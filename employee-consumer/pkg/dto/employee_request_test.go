package dto

import (
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestEmployee_ToString(t *testing.T) {
	var employee EmployeeRequest
	employee.FirstName = "Marcos"
	employee.LastName = "Luna"
	employee.SecondLastName = "Valdez"
	employee.DateOfBirth = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.DateOfEmployment = time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC)
	employee.Status = constants.ACTIVE

	jsonEmployee := `{"firstName":"Marcos","lastName":"Luna","secondLastName":"Valdez","dateOfBirth":"1994-04-25T08:00:00Z","dateOfEmployment":"1994-04-25T08:00:00Z","status":"ACTIVE"}`

	assert.Equal(t, jsonEmployee, employee.ToString())
}
