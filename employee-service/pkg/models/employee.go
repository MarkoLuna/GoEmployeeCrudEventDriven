package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MarkoLuna/EmployeeService/pkg/constants"
)

type Employee struct {
	Id               string                   `json:"id" swaggertype:"string" example:"b836ce65-76ab-42c8-b7b8-63ed432963c2"`
	FirstName        string                   `json:"firstName" validate:"required" swaggertype:"string" example:"Marcos"`
	LastName         string                   `json:"lastName" validate:"required" swaggertype:"string" example:"Luna"`
	SecondLastName   string                   `json:"secondLastName" validate:"required" swaggertype:"string" example:"Valdez"`
	DateOfBirth      time.Time                `json:"dateOfBirth" validate:"required" swaggertype:"string" example:"1994-04-25T12:00:00Z"`
	DateOfEmployment time.Time                `json:"dateOfEmployment" validate:"required" swaggertype:"string" example:"1994-04-25T12:00:00Z"`
	Status           constants.EmployeeStatus `json:"status" validate:"EmployeeStatusValid" swaggertype:"string" enums:"ACTIVE,INACTIVE"`
}

func (e Employee) ToString() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}
