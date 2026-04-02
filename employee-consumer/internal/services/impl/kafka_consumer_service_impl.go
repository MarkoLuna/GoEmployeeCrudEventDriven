package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/internal/dto"
	"github.com/MarkoLuna/EmployeeConsumer/internal/services"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO add concurrency
// TODO add policy retry

var (
	employeeInsertTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_INSERT_TOPIC", "employee-insert.v1")
	employeeUpdateTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPDATE_TOPIC", "employee-update.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

// KafkaConsumer defines the subset of kafka.Consumer methods used by this service.
type KafkaConsumer interface {
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Close() error
}

type KafkaConsumerServiceImpl struct {
	consumer        KafkaConsumer
	employeeService services.EmployeeService
}

func NewKafkaConsumerService(
	kafkaConsumer KafkaConsumer,
	employeeService services.EmployeeService) services.KafkaConsumerService {
	return &KafkaConsumerServiceImpl{
		consumer:        kafkaConsumer,
		employeeService: employeeService,
	}
}

func isConsumerEnabled() bool {
	enabled := utils.GetEnv("KAFKA_CONSUMER_ENABLED", "true")
	log.Printf("Consumer enabled: %s", enabled)
	consumers_enabled, _ := strconv.ParseBool(enabled)
	return consumers_enabled
}

func (kSrv *KafkaConsumerServiceImpl) Listen() error {
	if !isConsumerEnabled() {
		return nil
	}

	topics := []string{employeeInsertTopic, employeeUpdateTopic, employeeDeleteTopic}
	log.Printf("Listening for employee events on topics: %v", topics)

	err := kSrv.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	for {
		msg, err := kSrv.consumer.ReadMessage(-1)
		if err != nil {
			return fmt.Errorf("error reading message: %w", err)
		}

		if msg == nil {
			continue
		}

		log.Printf("Received message from topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))

		switch *msg.TopicPartition.Topic {
		case employeeInsertTopic:
			kSrv.handleInsert(msg)
		case employeeUpdateTopic:
			kSrv.handleUpdate(msg)
		case employeeDeleteTopic:
			kSrv.handleDelete(msg)
		default:
			log.Printf("received message from unknown topic: %s", *msg.TopicPartition.Topic)
		}
	}
}

func (kSrv *KafkaConsumerServiceImpl) handleInsert(msg *kafka.Message) {
	var employeeMessage dto.EmployeeMessage
	err := json.Unmarshal(msg.Value, &employeeMessage)
	if err != nil {
		log.Printf("Error decoding insert message: %v", err)
		return
	}

	fmt.Printf("Received Employee for creation: %v\n", employeeMessage)
	created, err := kSrv.employeeService.CreateEmployee(employeeMessage.EmployeeInfo)
	if err != nil {
		log.Printf("error creating employee %v: %v", employeeMessage, err)
		return
	}

	log.Printf("employee created successfully with id: %s", created.Id)
}

func (kSrv *KafkaConsumerServiceImpl) handleUpdate(msg *kafka.Message) {
	var employeeMessage dto.EmployeeMessage
	err := json.Unmarshal(msg.Value, &employeeMessage)
	if err != nil {
		log.Printf("Error decoding update message: %v", err)
		return
	}

	fmt.Printf("Received Employee for update: %v\n", employeeMessage)
	updated, err := kSrv.employeeService.UpdateEmployee(employeeMessage.ID, employeeMessage.EmployeeInfo)
	if err != nil {
		log.Printf("error updating employee %v: %v", employeeMessage, err)
		return
	}

	log.Printf("employee updated successfully with id: %s", updated.Id)
}

func (kSrv *KafkaConsumerServiceImpl) handleDelete(msg *kafka.Message) {
	employeeId := string(msg.Value)
	fmt.Printf("Received Employee for deletion: %s\n", employeeId)

	err := kSrv.employeeService.DeleteEmployeeById(employeeId)
	if err != nil {
		log.Printf("error deleting [employeeId: %s] error: %s", employeeId, err)
		return
	}

	log.Printf("employee deleted successfully with id: %s", employeeId)
}
