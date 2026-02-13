package impl

import (
	"os"
	"testing"

	appConfig "github.com/MarkoLuna/EmployeeConsumer/pkg/app/config"
	"github.com/stretchr/testify/assert"
)

func TestKafkaConsumerService_ListenEmployeeUpsert_NoErrorsWhenConsumerIsDisabled(t *testing.T) {

	consumer, _ := appConfig.NewKafkaConsumer()
	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	kafkaConsumerService := NewKafkaConsumerService(consumer)

	err := kafkaConsumerService.ListenEmployeeUpsert()

	assert.NoError(t, err)
}

func TestKafkaConsumerService_ListenEmployeeUpsert_ErrorsWhenConsumerIsEnabledAndKafkaFails(t *testing.T) {

	consumer, _ := appConfig.NewKafkaConsumer()
	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	kafkaConsumerService := NewKafkaConsumerService(consumer)

	err := kafkaConsumerService.ListenEmployeeUpsert()

	assert.Error(t, err)
}

func TestKafkaConsumerService_ListenEmployeeDeletio_NoErrorsWhenConsumerIsDisabled(t *testing.T) {

	consumer, _ := appConfig.NewKafkaConsumer()
	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	kafkaConsumerService := NewKafkaConsumerService(consumer)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.NoError(t, err)
}

func TestKafkaConsumerService_ListenEmployeeDeletio_ErrorsWhenConsumerIsEnabledAndKafkaFails(t *testing.T) {

	consumer, _ := appConfig.NewKafkaConsumer()
	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	kafkaConsumerService := NewKafkaConsumerService(consumer)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.Error(t, err)
}
