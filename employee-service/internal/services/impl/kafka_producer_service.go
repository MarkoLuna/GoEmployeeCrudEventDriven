package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	producerDeliveryTimeout = 10 * time.Second
)

var (
	employeeInsertTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_INSERT_TOPIC", "employee-insert.v1")
	employeeUpdateTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_UPDATE_TOPIC", "employee-update.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

type KafkaProducerService struct {
	producer *kafka.Producer
}

func NewKafkaProducerService(kafkaProducer *kafka.Producer) KafkaProducerService {
	return KafkaProducerService{producer: kafkaProducer}
}

func (kSrv KafkaProducerService) produceWithDelivery(topic string, value []byte) error {
	deliveryChan := make(chan kafka.Event, 1)

	err := kSrv.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny},
		Value: value,
	}, deliveryChan)

	if err != nil {
		log.Printf("error producing message: %v for topic: %s", err, topic)
		return err
	}

	select {
	case event := <-deliveryChan:
		msg, ok := event.(*kafka.Message)
		if !ok {
			return fmt.Errorf("unexpected delivery event type for topic: %s", topic)
		}
		if msg.TopicPartition.Error != nil {
			log.Printf("delivery failed for topic %s: %v", topic, msg.TopicPartition.Error)
			return msg.TopicPartition.Error
		}
		log.Printf("message delivered to topic: %s [%d] at offset %v",
			topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
		return nil

	case <-time.After(producerDeliveryTimeout):
		return fmt.Errorf("delivery timeout for topic: %s", topic)
	}
}

func (kSrv KafkaProducerService) SendInsert(employee dto.EmployeeMessage) error {
	value, err := json.Marshal(employee)
	if err != nil {
		return fmt.Errorf("failed to marshal employee: %w", err)
	}

	return kSrv.produceWithDelivery(employeeInsertTopic, value)
}

func (kSrv KafkaProducerService) SendUpdate(employee dto.EmployeeMessage) error {
	value, err := json.Marshal(employee)
	if err != nil {
		return fmt.Errorf("failed to marshal employee: %w", err)
	}

	return kSrv.produceWithDelivery(employeeUpdateTopic, value)
}

func (kSrv KafkaProducerService) SendDelete(employeeId string) error {
	return kSrv.produceWithDelivery(employeeDeleteTopic, []byte(employeeId))
}
