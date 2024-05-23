package impl

import (
	"testing"

	"github.com/MarkoLuna/EmployeeService/pkg/dto"
	"github.com/stretchr/testify/assert"
)

func TestKafkaProducerService_SendUpsert(t *testing.T) {

	producer, _ := BuildKafkaProducer()
	kafkaProducerService := NewKafkaProducerService(producer)

	err := kafkaProducerService.SendUpsert(dto.EmployeeRequest{})

	assert.NoError(t, err)
}

func TestKafkaProducerService_SendDelete(t *testing.T) {

	producer, _ := BuildKafkaProducer()
	kafkaProducerService := NewKafkaProducerService(producer)

	err := kafkaProducerService.SendDelete(dto.EmployeeRequest{})

	assert.NoError(t, err)
}
