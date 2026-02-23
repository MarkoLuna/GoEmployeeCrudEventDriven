package services

import (
	"github.com/MarkoLuna/EmployeeService/pkg/dto"
)

type KafkaProducerService interface {
	SendDelete(employee string) error
	SendInsert(employee dto.EmployeeMessage) error
	SendUpdate(employee dto.EmployeeMessage) error
}
