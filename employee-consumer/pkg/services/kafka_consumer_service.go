package services

type KafkaConsumerService interface {
	ListenEmployeeInsert()
	ListenEmployeeUpdate()
	ListenEmployeeDeletion()
}
