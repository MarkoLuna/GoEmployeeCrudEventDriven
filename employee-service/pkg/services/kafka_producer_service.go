package services

import (
	"github.com/MarkoLuna/EmployeeService/pkg/dto"
)

type KafkaProducerService interface {
	SendDelete(employee string) error
	SendUpsert(employee dto.EmployeeRequest) error
}
