package impl

import (
	"testing"

	"github.com/MarkoLuna/EmployeeService/internal/config"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/stretchr/testify/assert"
)

func TestKafkaProducerService_SendInsert(t *testing.T) {

	producer, _ := config.NewKafkaProducer()
	kafkaProducerService := NewKafkaProducerService(producer)

	err := kafkaProducerService.SendInsert(dto.EmployeeMessage{})

	assert.NoError(t, err)
}

func TestKafkaProducerService_SendUpdate(t *testing.T) {

	producer, _ := config.NewKafkaProducer()
	kafkaProducerService := NewKafkaProducerService(producer)

	err := kafkaProducerService.SendUpdate(dto.EmployeeMessage{})

	assert.NoError(t, err)
}

func TestKafkaProducerService_SendDelete(t *testing.T) {

	producer, _ := config.NewKafkaProducer()
	kafkaProducerService := NewKafkaProducerService(producer)

	err := kafkaProducerService.SendDelete("employeeId")

	assert.NoError(t, err)
}
