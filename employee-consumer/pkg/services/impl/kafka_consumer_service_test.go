package impl

import (
	"os"
	"testing"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/repositories"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestKafkaConsumerService_ListenEmployeeUpsert_NoErrorsWhenConsumerIsDisabled(t *testing.T) {

	consumer, _ := BuildKafkaConsumer()
	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := services.NewEmployeeService(employeeRepository)
	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	kafkaConsumerService := NewKafkaConsumerService(consumer, employeeService)

	err := kafkaConsumerService.ListenEmployeeUpsert()

	assert.NoError(t, err)
}

func TestKafkaConsumerService_ListenEmployeeUpsert_ErrorsWhenConsumerIsEnabledAndKafkaFails(t *testing.T) {

	consumer, _ := BuildKafkaConsumer()
	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := services.NewEmployeeService(employeeRepository)
	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	kafkaConsumerService := NewKafkaConsumerService(consumer, employeeService)

	err := kafkaConsumerService.ListenEmployeeUpsert()

	assert.Error(t, err)
}

func TestKafkaConsumerService_ListenEmployeeDeletio_NoErrorsWhenConsumerIsDisabled(t *testing.T) {

	consumer, _ := BuildKafkaConsumer()
	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := services.NewEmployeeService(employeeRepository)
	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	kafkaConsumerService := NewKafkaConsumerService(consumer, employeeService)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.NoError(t, err)
}

func TestKafkaConsumerService_ListenEmployeeDeletio_ErrorsWhenConsumerIsEnabledAndKafkaFails(t *testing.T) {

	consumer, _ := BuildKafkaConsumer()
	employeeRepository := repositories.NewEmployeeRepositoryStub()
	employeeService := services.NewEmployeeService(employeeRepository)
	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	kafkaConsumerService := NewKafkaConsumerService(consumer, employeeService)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.Error(t, err)
}
