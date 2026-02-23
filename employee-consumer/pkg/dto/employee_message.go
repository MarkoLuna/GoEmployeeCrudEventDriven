package dto

import (
	"encoding/json"
	"log"
)

type EmployeeMessage struct {
	ID           string
	EmployeeInfo EmployeeRequest
}

func (e EmployeeMessage) ToString() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}
