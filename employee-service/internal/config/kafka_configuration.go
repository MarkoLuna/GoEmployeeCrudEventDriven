package config

import (
	"github.com/MarkoLuna/EmployeeService/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	bootstrapServers = utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
)

func NewKafkaProducer() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})

	return p, err
}
