package services

import (
	"github.com/MarkoLuna/EmployeeService/pkg/dto"
)

type NewKafkaProducerService interface {
	SendDelete(employee dto.EmployeeRequest) error
	SendUpsert(employee dto.EmployeeRequest) error
}
