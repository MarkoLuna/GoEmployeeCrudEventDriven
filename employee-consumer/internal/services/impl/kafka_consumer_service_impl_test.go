package impl

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/internal/dto"
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

func TestKafkaConsumerService_Listen_NoErrorsWhenConsumerIsDisabled(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.Listen()

	assert.NoError(t, err)
	mockConsumer.AssertNotCalled(t, "SubscribeTopics")
}

func TestKafkaConsumerService_Listen_SuccessfullyDispatchesMessages(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

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
		mockEmployeeService,
	)

	err := kafkaConsumerService.Listen()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stop loop")

	mockEmployeeService.AssertCalled(t, "CreateEmployee", employeeRequest)
	mockEmployeeService.AssertCalled(t, "UpdateEmployee", "123", employeeRequest)
	mockEmployeeService.AssertCalled(t, "DeleteEmployeeById", deleteId)
}

func TestKafkaConsumerService_Listen_ContinuesOnIndividualErrors(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	topicInsert := "employee-insert.v1"
	msgInvalid := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicInsert},
		Value:          []byte("invalid-json"),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(msgInvalid, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.Listen()
	assert.Error(t, err)

	mockEmployeeService.AssertNotCalled(t, "CreateEmployee")
}
