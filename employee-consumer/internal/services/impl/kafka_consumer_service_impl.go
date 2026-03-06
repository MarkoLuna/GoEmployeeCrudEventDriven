package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/internal/dto"
	"github.com/MarkoLuna/EmployeeConsumer/internal/services"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO add concurrency
// TODO add policy retry

var (
	employeeInsertTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPSERT_TOPIC", "employee-insert.v1")
	employeeUpdateTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPSERT_TOPIC", "employee-update.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

// KafkaConsumer defines the subset of kafka.Consumer methods used by this service.
type KafkaConsumer interface {
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Close() error
}

type KafkaConsumerServiceImpl struct {
	consumer        KafkaConsumer
	employeeService services.EmployeeService
}

func NewKafkaConsumerService(
	kafkaConsumer KafkaConsumer,
	employeeService services.EmployeeService) services.KafkaConsumerService {
	return KafkaConsumerServiceImpl{
		consumer:        kafkaConsumer,
		employeeService: employeeService,
	}
}

func isConsumerEnabled() bool {
	enabled := utils.GetEnv("KAFKA_CONSUMER_ENABLED", "true")
	log.Printf("Consumer enabled: %s", enabled)
	consumers_enabled, _ := strconv.ParseBool(enabled)
	return consumers_enabled
}

func (kSrv KafkaConsumerServiceImpl) ListenEmployeeInsert() error {

	if isConsumerEnabled() {
		log.Printf("Listening for employee insert on topic: %s", employeeInsertTopic)
		kSrv.consumer.SubscribeTopics([]string{employeeInsertTopic}, nil)

		for {
			msg, err := kSrv.consumer.ReadMessage(time.Duration(-1))
			if err == nil {
				log.Printf("Received message: %v\n", msg)
				var employeeMessage dto.EmployeeMessage
				err := json.Unmarshal(msg.Value, &employeeMessage)
				if err != nil {
					fmt.Printf("Error decoding message: %v\n", err)
					continue
				}

				fmt.Printf("Received Employee for creation: %v\n", employeeMessage)
				created, err := kSrv.employeeService.CreateEmployee(employeeMessage.EmployeeInfo)
				if err != nil {
					log.Printf("error creating employee %v", employeeMessage)
					continue
				}

				log.Printf("employee created successfully with id: %s", created.Id)
			} else {
				log.Printf("error reading message: %v", err)
				continue
			}
		}
	}
	return nil
}

func (kSrv KafkaConsumerServiceImpl) ListenEmployeeUpdate() error {

	if isConsumerEnabled() {
		log.Printf("Listening for employee update on topic: %s", employeeUpdateTopic)
		kSrv.consumer.SubscribeTopics([]string{employeeUpdateTopic}, nil)

		for {
			msg, err := kSrv.consumer.ReadMessage(time.Duration(-1))
			if err == nil {
				log.Printf("Received message: %v\n", msg)
				var employeeMessage dto.EmployeeMessage
				err := json.Unmarshal(msg.Value, &employeeMessage)
				if err != nil {
					fmt.Printf("Error decoding message: %v\n", err)
					continue
				}

				fmt.Printf("Received Employee for update: %v\n", employeeMessage)
				updated, err := kSrv.employeeService.UpdateEmployee(employeeMessage.ID,
					employeeMessage.EmployeeInfo)
				if err != nil {
					log.Printf("error updating employee %v", employeeMessage)
					continue
				}

				log.Printf("employee updated successfully with id: %s", updated.Id)
				continue
			} else {
				log.Printf("error reading message: %v", err)
				continue
			}
		}
	}
	return nil
}

func (kSrv KafkaConsumerServiceImpl) ListenEmployeeDeletion() error {

	if isConsumerEnabled() {
		log.Printf("Listening for employee deletion on topic: %s", employeeDeleteTopic)
		kSrv.consumer.SubscribeTopics([]string{employeeDeleteTopic}, nil)

		for {
			msg, err := kSrv.consumer.ReadMessage(time.Duration(-1))
			if err == nil {
				log.Printf("Received message: %v\n", msg)
				var employeeId string = string(msg.Value)

				fmt.Printf("Received Employee for deletion: %s\n", employeeId)
				err := kSrv.employeeService.DeleteEmployeeById(employeeId)
				if err != nil {
					log.Printf("error deleting [employeeId: %s] error: %s", employeeId, err)
					continue
				}

				log.Printf("employee deleted successfully with id: %s", employeeId)
				continue
			} else {
				log.Printf("error reading message: %v", err)
				continue
			}
		}
	}
	return nil
}
