package services

type KafkaConsumerService interface {
	ListenEmployeeUpsert()
	ListenEmployeeDeletion()
}
