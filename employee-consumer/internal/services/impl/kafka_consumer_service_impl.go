package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/internal/dto"
	"github.com/MarkoLuna/EmployeeConsumer/internal/services"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO add policy retry

var (
	employeeInsertTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_INSERT_TOPIC", "employee-insert.v1")
	employeeUpdateTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPDATE_TOPIC", "employee-update.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)

// KafkaConsumer defines the subset of kafka.Consumer methods used by this service.
type KafkaConsumer interface {
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Close() error
}

// KafkaConsumerServiceImpl fans each Kafka topic out to its own dedicated
// worker pool so insert, update and delete workloads scale independently.
type KafkaConsumerServiceImpl struct {
	consumer           KafkaConsumer
	employeeService    services.EmployeeService
	insertWorkerCount  int
	updateWorkerCount  int
	deleteWorkerCount  int
}

func NewKafkaConsumerService(
	kafkaConsumer KafkaConsumer,
	employeeService services.EmployeeService,
) services.KafkaConsumerService {
	return &KafkaConsumerServiceImpl{
		consumer:          kafkaConsumer,
		employeeService:   employeeService,
		insertWorkerCount: utils.ParseIntEnv("KAFKA_CONSUMER_MAX_WORKERS_INSERT", 3),
		updateWorkerCount: utils.ParseIntEnv("KAFKA_CONSUMER_MAX_WORKERS_UPDATE", 3),
		deleteWorkerCount: utils.ParseIntEnv("KAFKA_CONSUMER_MAX_WORKERS_DELETE", 3),
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func isConsumerEnabled() bool {
	enabled := utils.ParseBoolEnv("KAFKA_CONSUMER_ENABLED", true)
	log.Printf("Consumer enabled: %v", enabled)
	return enabled
}

// startWorkerPool launches n goroutines that each drain ch by calling handler.
// Each worker logs its topic and ID on start/stop. The WaitGroup is decremented
// when a worker exits so the caller can wait for full drain after closing ch.
func startWorkerPool(
	wg *sync.WaitGroup,
	topic string,
	n int,
	ch <-chan *kafka.Message,
	handler func(*kafka.Message),
) {
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("[%s] worker %d started", topic, workerID)
			for msg := range ch {
				handler(msg)
			}
			log.Printf("[%s] worker %d stopped", topic, workerID)
		}(i)
	}
}

// ── Listen ───────────────────────────────────────────────────────────────────

func (kSrv *KafkaConsumerServiceImpl) Listen() error {
	if !isConsumerEnabled() {
		return nil
	}

	topics := []string{employeeInsertTopic, employeeUpdateTopic, employeeDeleteTopic}
	log.Printf(
		"Listening on topics: %v (insert workers: %d, update workers: %d, delete workers: %d)",
		topics,
		kSrv.insertWorkerCount, kSrv.updateWorkerCount, kSrv.deleteWorkerCount,
	)

	if err := kSrv.consumer.SubscribeTopics(topics, nil); err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	// Per-topic buffered channels — buffer size equals worker count so the
	// reader is never blocked longer than one in-flight message per worker.
	insertCh := make(chan *kafka.Message, kSrv.insertWorkerCount)
	updateCh := make(chan *kafka.Message, kSrv.updateWorkerCount)
	deleteCh := make(chan *kafka.Message, kSrv.deleteWorkerCount)

	var wg sync.WaitGroup
	startWorkerPool(&wg, employeeInsertTopic, kSrv.insertWorkerCount, insertCh, kSrv.handleInsert)
	startWorkerPool(&wg, employeeUpdateTopic, kSrv.updateWorkerCount, updateCh, kSrv.handleUpdate)
	startWorkerPool(&wg, employeeDeleteTopic, kSrv.deleteWorkerCount, deleteCh, kSrv.handleDelete)

	for {
		msg, err := kSrv.consumer.ReadMessage(-1)
		if err != nil {
			close(insertCh)
			close(updateCh)
			close(deleteCh)
			wg.Wait()
			return fmt.Errorf("error reading message: %w", err)
		}

		if msg == nil || msg.TopicPartition.Topic == nil {
			continue
		}

		switch *msg.TopicPartition.Topic {
		case employeeInsertTopic:
			insertCh <- msg
		case employeeUpdateTopic:
			updateCh <- msg
		case employeeDeleteTopic:
			deleteCh <- msg
		default:
			log.Printf("received message from unknown topic: %s", *msg.TopicPartition.Topic)
		}
	}
}

// ── Handlers ─────────────────────────────────────────────────────────────────

func (kSrv *KafkaConsumerServiceImpl) handleInsert(msg *kafka.Message) {
	var employeeMessage dto.EmployeeMessage
	err := json.Unmarshal(msg.Value, &employeeMessage)
	if err != nil {
		log.Printf("Error decoding insert message: %v", err)
		return
	}

	fmt.Printf("Received Employee for creation: %v\n", employeeMessage)
	created, err := kSrv.employeeService.CreateEmployee(employeeMessage.EmployeeInfo)
	if err != nil {
		log.Printf("error creating employee %v: %v", employeeMessage, err)
		return
	}

	log.Printf("employee created successfully with id: %s", created.Id)
}

func (kSrv *KafkaConsumerServiceImpl) handleUpdate(msg *kafka.Message) {
	var employeeMessage dto.EmployeeMessage
	err := json.Unmarshal(msg.Value, &employeeMessage)
	if err != nil {
		log.Printf("Error decoding update message: %v", err)
		return
	}

	fmt.Printf("Received Employee for update: %v\n", employeeMessage)
	updated, err := kSrv.employeeService.UpdateEmployee(employeeMessage.ID, employeeMessage.EmployeeInfo)
	if err != nil {
		log.Printf("error updating employee %v: %v", employeeMessage, err)
		return
	}

	log.Printf("employee updated successfully with id: %s", updated.Id)
}

func (kSrv *KafkaConsumerServiceImpl) handleDelete(msg *kafka.Message) {
	employeeId := string(msg.Value)
	fmt.Printf("Received Employee for deletion: %s\n", employeeId)

	err := kSrv.employeeService.DeleteEmployeeById(employeeId)
	if err != nil {
		log.Printf("error deleting [employeeId: %s] error: %s", employeeId, err)
		return
	}

	log.Printf("employee deleted successfully with id: %s", employeeId)
}
