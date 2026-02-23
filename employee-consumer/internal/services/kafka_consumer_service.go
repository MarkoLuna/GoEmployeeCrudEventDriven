package services

type KafkaConsumerService interface {
	ListenEmployeeInsert() error
	ListenEmployeeUpdate() error
	ListenEmployeeDeletion() error
}
