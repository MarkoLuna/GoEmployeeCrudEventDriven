package impl

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/services"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// TODO add concurrency
// TODO add policy retry

var (
	bootstrapServers    = utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	groupId             = utils.GetEnv("KAFKA_CONSUMER_GROUP_ID", "employee-group")
	employeeUpsertTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPSERT_TOPIC", "employee-upsert.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

type KafkaConsumerService struct {
	consumer        *kafka.Consumer
	employeeService services.EmployeeService
}

func BuildKafkaConsumer() (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	return c, err
}

func NewKafkaConsumerService(kafkaConsumer *kafka.Consumer,
	employeeService services.EmployeeService) KafkaConsumerService {
	return KafkaConsumerService{consumer: kafkaConsumer,
		employeeService: employeeService}
}

func isConsumerEnabled() bool {
	enabled := utils.GetEnv("KAFKA_CONSUMER_ENABLED", "false")
	consumers_enabled, _ := strconv.ParseBool(enabled)
	return consumers_enabled
}

func (kSrv KafkaConsumerService) ListenEmployeeUpsert() error {

	if isConsumerEnabled() {
		kSrv.consumer.SubscribeTopics([]string{employeeUpsertTopic}, nil)

		for {
			msg, err := kSrv.consumer.ReadMessage(-1)
			if err == nil {
				var employee dto.EmployeeDto
				err := json.Unmarshal(msg.Value, &employee)
				if err != nil {
					fmt.Printf("Error decoding message: %v\n", err)
					continue
				}

				fmt.Printf("Received Employee: %+v\n", employee)
				exists := kSrv.employeeService.ExistsEmployeeById(employee.Id)

				// TODO create mappers from employeeDto to employeeRequest
				newEmployee := dto.EmployeeRequest{}
				newEmployee.DateOfBirth = employee.DateOfBirth
				newEmployee.DateOfEmployment = employee.DateOfEmployment
				newEmployee.FirstName = employee.FirstName
				newEmployee.LastName = employee.LastName
				newEmployee.SecondLastName = employee.SecondLastName
				newEmployee.Status = employee.Status

				if !exists {
					kSrv.employeeService.CreateEmployee(newEmployee)
				} else {
					kSrv.employeeService.UpdateEmployee(employee.Id, newEmployee)
				}

			}
			return err
		}
	}
	return nil
}

func (kSrv KafkaConsumerService) ListenEmployeeDeletion() error {

	if isConsumerEnabled() {
		kSrv.consumer.SubscribeTopics([]string{employeeDeleteTopic}, nil)

		for {
			msg, err := kSrv.consumer.ReadMessage(-1)
			if err == nil {
				var employee dto.EmployeeDto
				err := json.Unmarshal(msg.Value, &employee)
				if err != nil {
					fmt.Printf("Error decoding message: %v\n", err)
					continue
				}

				fmt.Printf("Received Employee: %+v\n", employee)
				kSrv.employeeService.DeleteEmployeeById(employee.Id)
			}
			return err
		}
	}
	return nil
}
