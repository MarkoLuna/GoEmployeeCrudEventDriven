package stubs

import (
	"fmt"

	"github.com/MarkoLuna/EmployeeService/pkg/dto"
)

type KafkaProducerServiceStub struct{}

func NewKafkaProducerServiceStub() KafkaProducerServiceStub {
	return KafkaProducerServiceStub{}
}

func (kSrv KafkaProducerServiceStub) SendUpsert(employee dto.EmployeeRequest) error {

	fmt.Println("Upsert send")
	return nil
}

func (kSrv KafkaProducerServiceStub) SendDelete(employee dto.EmployeeRequest) error {

	fmt.Println("Delete send")
	return nil
}
