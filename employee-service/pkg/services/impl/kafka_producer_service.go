package impl

import (
	"encoding/json"

	"github.com/MarkoLuna/EmployeeService/pkg/dto"
	"github.com/MarkoLuna/EmployeeService/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	bootstrapServers    = utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	employeeUpsertTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_UPSERT_TOPIC", "employee-upsert.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

type KafkaProducerService struct {
	producer kafka.Producer
}

func NewKafkaProducerService() KafkaProducerService {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})

	if err != nil {
		panic(err)
	}
	defer p.Close()

	return KafkaProducerService{producer: *p}
}

func (kSrv KafkaProducerService) SendUpsert(employee dto.EmployeeRequest) error {

	value, err := json.Marshal(employee)
	if err != nil {
		panic(err)
	}

	err = kSrv.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &employeeUpsertTopic,
			Partition: kafka.PartitionAny},
		Value: value,
	}, nil)

	return err
	/*
		if err != nil {
			panic(err)
		}
	*/
}

func (kSrv KafkaProducerService) SendDelete(employee dto.EmployeeRequest) error {

	value, err := json.Marshal(employee)
	if err != nil {
		panic(err)
	}

	err = kSrv.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &employeeDeleteTopic,
			Partition: kafka.PartitionAny},
		Value: value,
	}, nil)

	return err
	/*
		if err != nil {
			panic(err)
		}
	*/
}
