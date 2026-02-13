package config

import (
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	bootstrapServers = utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	groupId          = utils.GetEnv("KAFKA_CONSUMER_GROUP_ID", "employee-group")
)

func NewKafkaConsumer() (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	return c, err
}
