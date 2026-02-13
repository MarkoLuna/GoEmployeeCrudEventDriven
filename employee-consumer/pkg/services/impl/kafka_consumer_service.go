package impl

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO add concurrency
// TODO add policy retry

var (
	employeeUpsertTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPSERT_TOPIC", "employee-upsert.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

type KafkaConsumerService struct {
	consumer *kafka.Consumer
}

func NewKafkaConsumerService(kafkaConsumer *kafka.Consumer) KafkaConsumerService {
	return KafkaConsumerService{consumer: kafkaConsumer}
}

func isConsumerEnabled() bool {
	enabled := utils.GetEnv("KAFKA_CONSUMER_ENABLED", "true")
	consumers_enabled, _ := strconv.ParseBool(enabled)
	return consumers_enabled
}

func (kSrv KafkaConsumerService) ListenEmployeeUpsert() error {

	if isConsumerEnabled() {
		kSrv.consumer.SubscribeTopics([]string{employeeUpsertTopic}, nil)

		for {
			msg, err := kSrv.consumer.ReadMessage(-1)
			if err == nil {
				var employee dto.EmployeeRequest
				err := json.Unmarshal(msg.Value, &employee)
				if err != nil {
					fmt.Printf("Error decoding message: %v\n", err)
					continue
				}

				fmt.Printf("Received Employee: %+v\n", employee)
				// TODO process upserts
				continue
			}
			return err
		}
	}
	return nil
}

func (kSrv KafkaConsumerService) ListenEmployeeDeletion() error {

	if isConsumerEnabled() {
		kSrv.consumer.SubscribeTopics([]string{employeeDeleteTopic}, nil)

		for {
			msg, err := kSrv.consumer.ReadMessage(-1)
			if err == nil {
				var employee string = string(msg.Value)

				fmt.Printf("Received Employee: %+v\n", employee)
				// TODO process deletes
				continue
			}
			return err
		}
	}
	return nil
}
