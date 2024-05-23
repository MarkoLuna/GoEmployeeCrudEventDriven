package impl

import (
	"encoding/json"
	"fmt"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO add concurrency
// TODO add policy retry

var (
	bootstrapServers    = utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	enabled             = utils.GetEnv("KAFKA_CONSUMER_ENABLED", "false")
	groupId             = utils.GetEnv("KAFKA_CONSUMER_GROUP_ID", "employee-group")
	employeeUpsertTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPSERT_TOPIC", "employee-upsert.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

type KafkaConsumerService struct {
	consumer kafka.Consumer
}

func NewKafkaConsumerService() KafkaConsumerService {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	return KafkaConsumerService{consumer: *c}
}

func (kSrv KafkaConsumerService) ListenEmployeeUpsert() {
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
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func (kSrv KafkaConsumerService) ListenEmployeeDeletion() {
	kSrv.consumer.SubscribeTopics([]string{employeeDeleteTopic}, nil)

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
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
