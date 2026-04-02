package dto

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/constants"
)

type EmployeeRequest struct {
	FirstName        string                   `json:"firstName" validate:"required" swaggertype:"string" example:"Marcos"`
	LastName         string                   `json:"lastName" validate:"required" swaggertype:"string" example:"Luna"`
	SecondLastName   string                   `json:"secondLastName" validate:"required" swaggertype:"string" example:"Valdez"`
	DateOfBirth      time.Time                `json:"dateOfBirth" validate:"required" swaggertype:"string" example:"1994-04-25T12:00:00Z"`
	DateOfEmployment time.Time                `json:"dateOfEmployment" validate:"required" swaggertype:"string" example:"1994-04-25T12:00:00Z"`
	Status           constants.EmployeeStatus `json:"status" validate:"EmployeeStatusValid" swaggertype:"string" enums:"ACTIVE,INACTIVE"`
}

func (e EmployeeRequest) ToString() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Println("Error marshaling EmployeeRequest:", err)
		return ""
	}
	return string(out)
}

type EmployeeMessage struct {
	ID           string
	EmployeeInfo EmployeeRequest
}

func (e EmployeeMessage) ToString() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Println("Error marshaling EmployeeMessage:", err)
		return ""
	}
	return string(out)
}

type GetEmployeeRequest struct {
	EmployeeId string `json:"employeeId" validate:"required" swaggertype:"string" example:"e26b200a-a8d0-11e9-a2a3-2a2ae2dbcce4"`
}

func (e GetEmployeeRequest) ToString() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}
