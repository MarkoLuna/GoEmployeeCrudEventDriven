package stubs

import (
	"fmt"
)

type KafkaConsumerServiceStub struct {
}

func NewKafkaConsumerServiceStub() KafkaConsumerServiceStub {
	return KafkaConsumerServiceStub{}
}

func (kSrv KafkaConsumerServiceStub) ListenEmployeeDeletion() {

	fmt.Println("listening employee deletion")
}

func (kSrv KafkaConsumerServiceStub) ListenEmployeeUpsert() {

	fmt.Println("listening employee upsert")
}
