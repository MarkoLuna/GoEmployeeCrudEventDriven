package impl

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/MarkoLuna/EmployeeConsumer/internal/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKafkaConsumer is a mock implementation of kafka.Consumer
type MockKafkaConsumer struct {
	mock.Mock
	ReadMessageFunc func(time.Duration) (*kafka.Message, error)
}

func (m *MockKafkaConsumer) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {
	args := m.Called(topics, rebalanceCb)
	return args.Error(0)
}

func (m *MockKafkaConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	if m.ReadMessageFunc != nil {
		return m.ReadMessageFunc(timeout)
	}
	args := m.Called(timeout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kafka.Message), args.Error(1)
}

func (m *MockKafkaConsumer) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockKafkaProducer is a mock implementation of kafka.Producer
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	args := m.Called(msg, deliveryChan)
	return args.Error(0)
}

func (m *MockKafkaProducer) Close() {
	m.Called()
}

// MockEmployeeService is a mock implementation of EmployeeService
type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) CreateEmployee(employeeRequest dto.EmployeeRequest) (*models.Employee, error) {
	args := m.Called(employeeRequest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeService) UpdateEmployee(employeeId string, employeeRequest dto.EmployeeRequest) (*models.Employee, error) {
	args := m.Called(employeeId, employeeRequest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeService) DeleteEmployeeById(employeeId string) error {
	args := m.Called(employeeId)
	return args.Error(0)
}

func (m *MockEmployeeService) GetEmployees() ([]models.Employee, error) {
	args := m.Called()
	return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *MockEmployeeService) GetEmployeeById(employeeId string) (models.Employee, error) {
	args := m.Called(employeeId)
	return args.Get(0).(models.Employee), args.Error(1)
}

func createEmployeeMessageBytes(t *testing.T, id string, employeeInfo dto.EmployeeRequest) []byte {
	msg := dto.EmployeeMessage{
		ID:           id,
		EmployeeInfo: employeeInfo,
	}
	bytes, err := json.Marshal(msg)
	assert.NoError(t, err)
	return bytes
}

// newMockConsumerWithClose returns a MockKafkaConsumer with Close() already stubbed.
// All tests must call Close() because Listen always closes the consumer on exit.
func newMockConsumerWithClose() *MockKafkaConsumer {
	m := new(MockKafkaConsumer)
	m.On("Close").Return(nil)
	return m
}

func TestKafkaConsumerService_Listen_NoErrorsWhenConsumerIsDisabled(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)
	mockProducer := new(MockKafkaProducer)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockProducer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.Listen(context.Background())

	assert.NoError(t, err)
	mockConsumer.AssertNotCalled(t, "SubscribeTopics")
}

func TestKafkaConsumerService_Listen_SuccessfullyDispatchesMessages(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)
	mockProducer := new(MockKafkaProducer)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Setenv("KAFKA_CONSUMER_ENABLED", "false")

	employeeRequest := dto.EmployeeRequest{FirstName: "John"}
	insertBytes := createEmployeeMessageBytes(t, "", employeeRequest)
	updateBytes := createEmployeeMessageBytes(t, "123", employeeRequest)
	deleteId := "123"

	topicInsert := "employee-insert.v1"
	topicUpdate := "employee-update.v1"
	topicDelete := "employee-deletion.v1"

	msgInsert := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicInsert},
		Value:          insertBytes,
	}
	msgUpdate := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicUpdate},
		Value:          updateBytes,
	}
	msgDelete := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicDelete},
		Value:          []byte(deleteId),
	}

	messages := []*kafka.Message{msgInsert, msgUpdate, msgDelete}
	msgIndex := 0

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.ReadMessageFunc = func(timeout time.Duration) (*kafka.Message, error) {
		if msgIndex < len(messages) {
			msg := messages[msgIndex]
			msgIndex++
			return msg, nil
		}
		return nil, errors.New("stop loop")
	}

	mockEmployeeService.On("CreateEmployee", employeeRequest).Return(&models.Employee{Id: "123"}, nil)
	mockEmployeeService.On("UpdateEmployee", "123", employeeRequest).Return(&models.Employee{Id: "123"}, nil)
	mockEmployeeService.On("DeleteEmployeeById", deleteId).Return(nil)

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockProducer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.Listen(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stop loop")

	mockEmployeeService.AssertCalled(t, "CreateEmployee", employeeRequest)
	mockEmployeeService.AssertCalled(t, "UpdateEmployee", "123", employeeRequest)
	mockEmployeeService.AssertCalled(t, "DeleteEmployeeById", deleteId)
}

func TestKafkaConsumerService_Listen_ContinuesOnIndividualErrors(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	topicInsert := "employee-insert.v1"
	msgInvalid := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicInsert},
		Value:          []byte("invalid-json"),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", mock.Anything).Return(msgInvalid, nil).Once()
	mockConsumer.On("ReadMessage", mock.Anything).Return(nil, errors.New("stop loop")).Once()

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		new(MockKafkaProducer),
		mockEmployeeService,
	)

	err := kafkaConsumerService.Listen(context.Background())
	assert.Error(t, err)

	mockEmployeeService.AssertNotCalled(t, "CreateEmployee")
}

func TestKafkaConsumerService_Listen_ConcurrentWorkers(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	// Each topic gets its own worker pool size.
	os.Setenv("KAFKA_CONSUMER_MAX_WORKERS_INSERT", "3")
	os.Setenv("KAFKA_CONSUMER_MAX_WORKERS_UPDATE", "2")
	os.Setenv("KAFKA_CONSUMER_MAX_WORKERS_DELETE", "1")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")
	defer os.Unsetenv("KAFKA_CONSUMER_MAX_WORKERS_INSERT")
	defer os.Unsetenv("KAFKA_CONSUMER_MAX_WORKERS_UPDATE")
	defer os.Unsetenv("KAFKA_CONSUMER_MAX_WORKERS_DELETE")

	employeeRequest := dto.EmployeeRequest{FirstName: "John"}

	topicInsert := "employee-insert.v1"
	topicUpdate := "employee-update.v1"
	topicDelete := "employee-deletion.v1"

	insertBytes := createEmployeeMessageBytes(t, "", employeeRequest)
	updateBytes := createEmployeeMessageBytes(t, "123", employeeRequest)

	// 6 insert, 4 update, 2 delete = 12 total messages.
	var allMessages []*kafka.Message
	for i := 0; i < 6; i++ {
		allMessages = append(allMessages, &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicInsert, Partition: 0, Offset: kafka.Offset(i)},
			Value:          insertBytes,
		})
	}
	for i := 0; i < 4; i++ {
		allMessages = append(allMessages, &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicUpdate, Partition: 0, Offset: kafka.Offset(i)},
			Value:          updateBytes,
		})
	}
	for i := 0; i < 2; i++ {
		allMessages = append(allMessages, &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicDelete, Partition: 0, Offset: kafka.Offset(i)},
			Value:          []byte("123"),
		})
	}

	var msgIndex int32

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.ReadMessageFunc = func(timeout time.Duration) (*kafka.Message, error) {
		idx := atomic.AddInt32(&msgIndex, 1) - 1
		if int(idx) < len(allMessages) {
			return allMessages[idx], nil
		}
		return nil, errors.New("stop loop")
	}

	var insertCount, updateCount, deleteCount int32

	mockEmployeeService.On("CreateEmployee", employeeRequest).
		Return(&models.Employee{Id: "123"}, nil).
		Run(func(args mock.Arguments) { atomic.AddInt32(&insertCount, 1) })

	mockEmployeeService.On("UpdateEmployee", "123", employeeRequest).
		Return(&models.Employee{Id: "123"}, nil).
		Run(func(args mock.Arguments) { atomic.AddInt32(&updateCount, 1) })

	mockEmployeeService.On("DeleteEmployeeById", "123").
		Return(nil).
		Run(func(args mock.Arguments) { atomic.AddInt32(&deleteCount, 1) })

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		new(MockKafkaProducer),
		mockEmployeeService,
	)

	err := kafkaConsumerService.Listen(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stop loop")

	assert.Equal(t, int32(6), atomic.LoadInt32(&insertCount), "all 6 insert messages should be processed")
	assert.Equal(t, int32(4), atomic.LoadInt32(&updateCount), "all 4 update messages should be processed")
	assert.Equal(t, int32(2), atomic.LoadInt32(&deleteCount), "all 2 delete messages should be processed")
}

func TestKafkaConsumerService_Listen_RetriesOnTransientError(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	os.Setenv("KAFKA_CONSUMER_MAX_RETRIES", "3")
	os.Setenv("KAFKA_CONSUMER_RETRY_INITIAL_BACKOFF_MS", "1") // Fast retry for tests
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")
	defer os.Unsetenv("KAFKA_CONSUMER_MAX_RETRIES")
	defer os.Unsetenv("KAFKA_CONSUMER_RETRY_INITIAL_BACKOFF_MS")

	employeeRequest := dto.EmployeeRequest{FirstName: "Retry"}
	topicInsert := "employee-insert.v1"
	msgInsert := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicInsert, Partition: 0, Offset: 1},
		Value:          createEmployeeMessageBytes(t, "", employeeRequest),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", mock.Anything).Return(msgInsert, nil).Once()
	mockConsumer.On("ReadMessage", mock.Anything).Return(nil, errors.New("stop loop")).Once()

	// Fail twice, succeed on third attempt
	mockEmployeeService.On("CreateEmployee", employeeRequest).Return(nil, errors.New("transient error")).Twice()
	mockEmployeeService.On("CreateEmployee", employeeRequest).Return(&models.Employee{Id: "retry-id"}, nil).Once()

	kafkaConsumerService := NewKafkaConsumerService(mockConsumer, new(MockKafkaProducer), mockEmployeeService)
	err := kafkaConsumerService.Listen(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stop loop")
	mockEmployeeService.AssertExpectations(t)
}

func TestKafkaConsumerService_Listen_SkipsDuplicateMessages(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	employeeRequest := dto.EmployeeRequest{FirstName: "Duplicate"}
	topicInsert := "employee-insert.v1"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicInsert, Partition: 0, Offset: 100},
		Value:          createEmployeeMessageBytes(t, "", employeeRequest),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	// Deliver the SAME message twice
	mockConsumer.On("ReadMessage", mock.Anything).Return(msg, nil).Twice()
	mockConsumer.On("ReadMessage", mock.Anything).Return(nil, errors.New("stop loop")).Once()

	// Should only be called ONCE
	mockEmployeeService.On("CreateEmployee", employeeRequest).Return(&models.Employee{Id: "id"}, nil).Once()

	kafkaConsumerService := NewKafkaConsumerService(mockConsumer, new(MockKafkaProducer), mockEmployeeService)
	err := kafkaConsumerService.Listen(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stop loop")
	mockEmployeeService.AssertExpectations(t)
}

func TestKafkaConsumerService_Listen_SendsToDLTAfterFinalFailure(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)
	mockProducer := new(MockKafkaProducer)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	os.Setenv("KAFKA_CONSUMER_MAX_RETRIES", "2") // 2 attempts total
	os.Setenv("KAFKA_CONSUMER_RETRY_INITIAL_BACKOFF_MS", "1")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")
	defer os.Unsetenv("KAFKA_CONSUMER_MAX_RETRIES")
	defer os.Unsetenv("KAFKA_CONSUMER_RETRY_INITIAL_BACKOFF_MS")

	employeeRequest := dto.EmployeeRequest{FirstName: "DLT"}
	topicInsert := "employee-insert.v1"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicInsert, Partition: 0, Offset: 200},
		Value:          createEmployeeMessageBytes(t, "", employeeRequest),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", mock.Anything).Return(msg, nil).Once()
	mockConsumer.On("ReadMessage", mock.Anything).Return(nil, errors.New("stop loop")).Once()

	// Always fail
	mockEmployeeService.On("CreateEmployee", employeeRequest).Return(nil, errors.New("permanent transient failure")).Times(2)

	// Verify DLT production
	mockProducer.On("Produce", mock.MatchedBy(func(m *kafka.Message) bool {
		return *m.TopicPartition.Topic == "employee-insert.v1.dlt" &&
			string(m.Value) == string(msg.Value)
	}), mock.Anything).Return(nil).Once()

	kafkaConsumerService := NewKafkaConsumerService(mockConsumer, mockProducer, mockEmployeeService)
	err := kafkaConsumerService.Listen(context.Background())

	assert.Error(t, err)
	mockEmployeeService.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

// TestKafkaConsumerService_Listen_GracefulShutdownViaContext verifies that
// cancelling the context causes Listen to return nil (not an error).
func TestKafkaConsumerService_Listen_GracefulShutdownViaContext(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	// ReadMessage always times out so Listen keeps polling.
	mockConsumer.ReadMessageFunc = func(timeout time.Duration) (*kafka.Message, error) {
		return nil, kafka.NewError(kafka.ErrTimedOut, "timed out", false)
	}

	ctx, cancel := context.WithCancel(context.Background())

	kafkaConsumerService := NewKafkaConsumerService(mockConsumer, new(MockKafkaProducer), mockEmployeeService)

	done := make(chan error, 1)
	go func() {
		done <- kafkaConsumerService.Listen(ctx)
	}()

	// Give Listen time to enter its poll loop, then cancel.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.NoError(t, err, "graceful shutdown should return nil")
	case <-time.After(3 * time.Second):
		t.Fatal("Listen did not return after context cancellation")
	}
}

// TestKafkaConsumerService_Stop_IsIdempotent verifies Stop can be called
// multiple times without panicking.
func TestKafkaConsumerService_Stop_IsIdempotent(t *testing.T) {
	mockConsumer := newMockConsumerWithClose()
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.ReadMessageFunc = func(timeout time.Duration) (*kafka.Message, error) {
		return nil, kafka.NewError(kafka.ErrTimedOut, "timed out", false)
	}

	svc := NewKafkaConsumerService(mockConsumer, new(MockKafkaProducer), mockEmployeeService)

	done := make(chan error, 1)
	go func() {
		done <- svc.Listen(context.Background())
	}()

	time.Sleep(50 * time.Millisecond)

	// Call Stop multiple times — must not panic.
	assert.NotPanics(t, func() {
		svc.Stop()
		svc.Stop()
		svc.Stop()
	})

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("Listen did not return after Stop()")
	}
}
