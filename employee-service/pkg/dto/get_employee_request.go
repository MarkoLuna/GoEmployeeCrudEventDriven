package dto

import (
	"encoding/json"
	"log"
)

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
