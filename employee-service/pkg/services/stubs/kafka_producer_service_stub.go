package stubs

import (
	"fmt"

	"github.com/MarkoLuna/EmployeeService/pkg/dto"
)

type KafkaProducerServiceStub struct {
	err error
}

func NewKafkaProducerServiceStub() KafkaProducerServiceStub {
	return KafkaProducerServiceStub{err: nil}
}

func NewKafkaProducerServiceStubFromError(err error) KafkaProducerServiceStub {
	return KafkaProducerServiceStub{err: err}
}

func (kSrv KafkaProducerServiceStub) SendUpsert(employee dto.EmployeeRequest) error {

	fmt.Println("Upsert send")
	return kSrv.err
}

func (kSrv KafkaProducerServiceStub) SendDelete(employee string) error {

	fmt.Println("Delete send")
	return kSrv.err
}
