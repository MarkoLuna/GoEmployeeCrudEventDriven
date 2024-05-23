package dto

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/constants"
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
		log.Fatal(err)
	}

	return string(out)
}
