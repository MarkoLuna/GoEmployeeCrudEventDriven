package impl

import (
	"encoding/json"

	"github.com/MarkoLuna/EmployeeService/pkg/dto"
	"github.com/MarkoLuna/EmployeeService/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	employeeInsertTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_UPSERT_TOPIC", "employee-insert.v1")
	employeeUpdateTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_UPSERT_TOPIC", "employee-update.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

type KafkaProducerService struct {
	producer *kafka.Producer
}

func NewKafkaProducerService(kafkaProducer *kafka.Producer) KafkaProducerService {
	return KafkaProducerService{producer: kafkaProducer}
}

func (kSrv KafkaProducerService) SendInsert(employee dto.EmployeeMessage) error {

	value, err := json.Marshal(employee)
	if err != nil {
		panic(err)
	}

	err = kSrv.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &employeeInsertTopic,
			Partition: kafka.PartitionAny},
		Value: value,
	}, nil)

	return err
}

func (kSrv KafkaProducerService) SendUpdate(employee dto.EmployeeMessage) error {

	value, err := json.Marshal(employee)
	if err != nil {
		panic(err)
	}

	err = kSrv.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &employeeUpdateTopic,
			Partition: kafka.PartitionAny},
		Value: value,
	}, nil)

	return err
}

func (kSrv KafkaProducerService) SendDelete(employeeId string) error {

	err := kSrv.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &employeeDeleteTopic,
			Partition: kafka.PartitionAny},
		Value: []byte(employeeId),
	}, nil)

	return err
}
