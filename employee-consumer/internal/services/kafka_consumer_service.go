package services

type KafkaConsumerService interface {
	Listen() error
}
