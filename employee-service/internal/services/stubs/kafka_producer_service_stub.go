package stubs

import (
	"fmt"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
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

func (kSrv KafkaProducerServiceStub) SendInsert(employee dto.EmployeeMessage) error {

	fmt.Println("Insert send")
	return kSrv.err
}

func (kSrv KafkaProducerServiceStub) SendUpdate(employee dto.EmployeeMessage) error {

	fmt.Println("Update send")
	return kSrv.err
}

func (kSrv KafkaProducerServiceStub) SendDelete(employee string) error {

	fmt.Println("Delete send")
	return kSrv.err
}
