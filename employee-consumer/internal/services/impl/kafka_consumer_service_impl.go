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

var (
	employeeInsertTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_INSERT_TOPIC", "employee-insert.v1")
	employeeUpdateTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPDATE_TOPIC", "employee-update.v1")
	employeeDeleteTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")

	employeeInsertDLTTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_INSERT_DLT", "employee-insert.v1.dlt")
	employeeUpdateDLTTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_UPDATE_DLT", "employee-update.v1.dlt")
	employeeDeleteDLTTopic = utils.GetEnv("KAFKA_CONSUMER_EMPLOYEE_DELETE_DLT", "employee-deletion.v1.dlt")
)

type KafkaConsumer interface {
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Close() error
}

// KafkaProducer defines the subset of kafka.Producer methods used by this service.
type KafkaProducer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	Close()
}

// KafkaConsumerServiceImpl fans each Kafka topic out to its own dedicated
// worker pool so insert, update and delete workloads scale independently.
type KafkaConsumerServiceImpl struct {
	consumer            KafkaConsumer
	producer            KafkaProducer
	employeeService     services.EmployeeService
	insertWorkerCount   int
	updateWorkerCount   int
	deleteWorkerCount   int
	retryMaxAttempts    int
	retryInitialBackoff time.Duration
	processedKeys       sync.Map // key: "topic:partition:offset" → struct{}
}

func NewKafkaConsumerService(
	kafkaConsumer KafkaConsumer,
	kafkaProducer KafkaProducer,
	employeeService services.EmployeeService,
) services.KafkaConsumerService {
	return &KafkaConsumerServiceImpl{
		consumer:            kafkaConsumer,
		producer:            kafkaProducer,
		employeeService:     employeeService,
		insertWorkerCount:   utils.ParseIntEnv("KAFKA_CONSUMER_MAX_WORKERS_INSERT", 3),
		updateWorkerCount:   utils.ParseIntEnv("KAFKA_CONSUMER_MAX_WORKERS_UPDATE", 3),
		deleteWorkerCount:   utils.ParseIntEnv("KAFKA_CONSUMER_MAX_WORKERS_DELETE", 3),
		retryMaxAttempts:    utils.ParseIntEnv("KAFKA_CONSUMER_MAX_RETRIES", 3),
		retryInitialBackoff: time.Duration(utils.ParseIntEnv("KAFKA_CONSUMER_RETRY_INITIAL_BACKOFF_MS", 500)) * time.Millisecond,
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func isConsumerEnabled() bool {
	enabled := utils.ParseBoolEnv("KAFKA_CONSUMER_ENABLED", true)
	log.Printf("Consumer enabled: %v", enabled)
	return enabled
}

func messageKey(msg *kafka.Message) string {
	topic := ""
	if msg.TopicPartition.Topic != nil {
		topic = *msg.TopicPartition.Topic
	}
	return fmt.Sprintf("%s:%d:%d", topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
}

func (kSrv *KafkaConsumerServiceImpl) withRetry(fn func() error) error {
	maxAttempts := kSrv.retryMaxAttempts
	backoff := kSrv.retryInitialBackoff

	for i := 0; i < maxAttempts; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		if i == maxAttempts-1 {
			return err
		}
		log.Printf("Retry attempt %d/%d failed: %v. Retrying in %v", i+1, maxAttempts, err, backoff)
		time.Sleep(backoff)
		backoff *= 2
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}
	}
	return nil
}

func (kSrv *KafkaConsumerServiceImpl) produceToDLT(dltTopic string, originalMsg *kafka.Message, finalErr error) {
	if kSrv.producer == nil {
		log.Printf("DLT producer is not configured, cannot send to %s", dltTopic)
		return
	}

	originalTopic := ""
	if originalMsg.TopicPartition.Topic != nil {
		originalTopic = *originalMsg.TopicPartition.Topic
	}

	dltMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &dltTopic, Partition: kafka.PartitionAny},
		Value:          originalMsg.Value,
		Key:            originalMsg.Key,
		Headers: append(originalMsg.Headers,
			kafka.Header{Key: "x-retries", Value: []byte(fmt.Sprintf("%d", kSrv.retryMaxAttempts))},
			kafka.Header{Key: "x-error-message", Value: []byte(finalErr.Error())},
			kafka.Header{Key: "x-original-topic", Value: []byte(originalTopic)},
		),
	}

	err := kSrv.producer.Produce(dltMsg, nil)
	if err != nil {
		log.Printf("Failed to produce to DLT %s: %v", dltTopic, err)
		return
	}
	log.Printf("Successfully sent message to DLT %s after %d retries", dltTopic, kSrv.retryMaxAttempts)
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
	key := messageKey(msg)
	if _, loaded := kSrv.processedKeys.LoadOrStore(key, struct{}{}); loaded {
		log.Printf("Message with key %s already being processed or already processed, skipping", key)
		return
	}

	var employeeMessage dto.EmployeeMessage
	err := json.Unmarshal(msg.Value, &employeeMessage)
	if err != nil {
		log.Printf("Error decoding insert message (not retried): %v", err)
		kSrv.processedKeys.Delete(key) // Unblock after permanent error
		return
	}

	fmt.Printf("Received Employee for creation: %v\n", employeeMessage)
	err = kSrv.withRetry(func() error {
		created, err := kSrv.employeeService.CreateEmployee(employeeMessage.EmployeeInfo)
		if err != nil {
			return err
		}
		log.Printf("employee created successfully with id: %s", created.Id)
		return nil
	})

	if err != nil {
		log.Printf("Final failure creating employee after retries: %v", err)
		kSrv.produceToDLT(employeeInsertDLTTopic, msg, err)
		kSrv.processedKeys.Delete(key) // Unblock after final transient failure
	}
}

func (kSrv *KafkaConsumerServiceImpl) handleUpdate(msg *kafka.Message) {
	key := messageKey(msg)
	if _, loaded := kSrv.processedKeys.LoadOrStore(key, struct{}{}); loaded {
		log.Printf("Message with key %s already being processed or already processed, skipping", key)
		return
	}

	var employeeMessage dto.EmployeeMessage
	err := json.Unmarshal(msg.Value, &employeeMessage)
	if err != nil {
		log.Printf("Error decoding update message (not retried): %v", err)
		kSrv.processedKeys.Delete(key) // Unblock after permanent error
		return
	}

	fmt.Printf("Received Employee for update: %v\n", employeeMessage)
	err = kSrv.withRetry(func() error {
		updated, err := kSrv.employeeService.UpdateEmployee(employeeMessage.ID, employeeMessage.EmployeeInfo)
		if err != nil {
			return err
		}
		log.Printf("employee updated successfully with id: %s", updated.Id)
		return nil
	})

	if err != nil {
		log.Printf("Final failure updating employee after retries: %v", err)
		kSrv.produceToDLT(employeeUpdateDLTTopic, msg, err)
		kSrv.processedKeys.Delete(key) // Unblock after final transient failure
	}
}

func (kSrv *KafkaConsumerServiceImpl) handleDelete(msg *kafka.Message) {
	key := messageKey(msg)
	if _, loaded := kSrv.processedKeys.LoadOrStore(key, struct{}{}); loaded {
		log.Printf("Message with key %s already being processed or already processed, skipping", key)
		return
	}

	employeeId := string(msg.Value)
	fmt.Printf("Received Employee for deletion: %s\n", employeeId)

	err := kSrv.withRetry(func() error {
		err := kSrv.employeeService.DeleteEmployeeById(employeeId)
		if err != nil {
			return err
		}
		log.Printf("employee deleted successfully with id: %s", employeeId)
		return nil
	})

	if err != nil {
		log.Printf("Final failure deleting employee after retries: %v", err)
		kSrv.produceToDLT(employeeDeleteDLTTopic, msg, err)
		kSrv.processedKeys.Delete(key) // Unblock after final transient failure
	}
}
