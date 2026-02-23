package impl

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKafkaConsumer is a mock implementation of kafka.Consumer
type MockKafkaConsumer struct {
	mock.Mock
}

func (m *MockKafkaConsumer) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {
	args := m.Called(topics, rebalanceCb)
	return args.Error(0)
}

func (m *MockKafkaConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
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

func TestKafkaConsumerService_ListenEmployeeInsert_NoErrorsWhenConsumerIsDisabled(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeInsert()

	assert.NoError(t, err)
	mockConsumer.AssertNotCalled(t, "SubscribeTopics")
	mockEmployeeService.AssertNotCalled(t, "CreateEmployee")
}

func TestKafkaConsumerService_ListenEmployeeInsert_ReturnsErrorWhenReadMessageFails(t *testing.T) {
	// t.Skip("skiped until mockConsumer matches with expected type")
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	expectedError := errors.New("kafka read error")
	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, expectedError)

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeInsert()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockConsumer.AssertCalled(t, "SubscribeTopics", mock.Anything, mock.Anything)
	mockEmployeeService.AssertNotCalled(t, "CreateEmployee")
}

func TestKafkaConsumerService_ListenEmployeeInsert_SuccessfullyCreatesEmployee(t *testing.T) {
	t.Skip("skiped until mockConsumer matches with expected type")
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	employeeRequest := dto.EmployeeRequest{
		FirstName:        "John",
		LastName:         "Doe",
		SecondLastName:   "Smith",
		DateOfBirth:      time.Date(1990, time.January, 01, 8, 0, 0, 0, time.UTC),
		DateOfEmployment: time.Date(2000, time.January, 01, 8, 0, 0, 0, time.UTC),
		Status:           "Active",
	}

	messageBytes := createEmployeeMessageBytes(t, "", employeeRequest)
	kafkaMessage := &kafka.Message{
		Value: messageBytes,
	}

	createdEmployee := &models.Employee{
		Id:               "123",
		FirstName:        "John",
		LastName:         "Doe",
		SecondLastName:   "Smith",
		DateOfBirth:      time.Date(1990, time.January, 01, 8, 0, 0, 0, time.UTC),
		DateOfEmployment: time.Date(2000, time.January, 01, 8, 0, 0, 0, time.UTC),
		Status:           "Active",
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(kafkaMessage, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()
	mockEmployeeService.On("CreateEmployee", employeeRequest).Return(createdEmployee, nil)

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeInsert()

	assert.Error(t, err) // Will error on second ReadMessage to stop the loop
	mockConsumer.AssertCalled(t, "SubscribeTopics", mock.Anything, mock.Anything)
	mockEmployeeService.AssertCalled(t, "CreateEmployee", employeeRequest)
}

func TestKafkaConsumerService_ListenEmployeeInsert_ContinuesOnInvalidJSON(t *testing.T) {
	t.Skip("skiped until mockConsumer matches with expected type")
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaMessage := &kafka.Message{
		Value: []byte("invalid json"),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(kafkaMessage, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeInsert()

	assert.Error(t, err) // Will error on second ReadMessage to stop the loop
	mockConsumer.AssertCalled(t, "SubscribeTopics", mock.Anything, mock.Anything)
	mockEmployeeService.AssertNotCalled(t, "CreateEmployee")
}

func TestKafkaConsumerService_ListenEmployeeInsert_ContinuesOnServiceError(t *testing.T) {
	t.Skip("skiped until mockConsumer matches with expected type")
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	employeeRequest := dto.EmployeeRequest{
		FirstName:        "John",
		LastName:         "Doe",
		SecondLastName:   "Smith",
		DateOfBirth:      time.Date(1990, time.January, 01, 8, 0, 0, 0, time.UTC),
		DateOfEmployment: time.Date(2000, time.January, 01, 8, 0, 0, 0, time.UTC),
		Status:           "Active",
	}

	messageBytes := createEmployeeMessageBytes(t, "", employeeRequest)
	kafkaMessage := &kafka.Message{
		Value: messageBytes,
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(kafkaMessage, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()
	mockEmployeeService.On("CreateEmployee", employeeRequest).Return(nil, errors.New("service error"))

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeInsert()

	assert.Error(t, err) // Will error on second ReadMessage to stop the loop
	mockEmployeeService.AssertCalled(t, "CreateEmployee", employeeRequest)
}

func TestKafkaConsumerService_ListenEmployeeUpdate_NoErrorsWhenConsumerIsDisabled(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeUpdate()

	assert.NoError(t, err)
	mockConsumer.AssertNotCalled(t, "SubscribeTopics")
	mockEmployeeService.AssertNotCalled(t, "UpdateEmployee")
}

func TestKafkaConsumerService_ListenEmployeeUpdate_ReturnsErrorWhenReadMessageFails(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	expectedError := errors.New("kafka read error")
	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, expectedError)

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeUpdate()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockConsumer.AssertCalled(t, "SubscribeTopics", mock.Anything, mock.Anything)
	mockEmployeeService.AssertNotCalled(t, "UpdateEmployee")
}

func TestKafkaConsumerService_ListenEmployeeUpdate_SuccessfullyUpdatesEmployee(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	employeeId := "123"
	employeeRequest := dto.EmployeeRequest{
		FirstName:        "John",
		LastName:         "Doe",
		SecondLastName:   "Smith",
		DateOfBirth:      time.Date(1990, time.January, 01, 8, 0, 0, 0, time.UTC),
		DateOfEmployment: time.Date(2000, time.January, 01, 8, 0, 0, 0, time.UTC),
		Status:           "Active",
	}

	messageBytes := createEmployeeMessageBytes(t, employeeId, employeeRequest)
	kafkaMessage := &kafka.Message{
		Value: messageBytes,
	}

	updatedEmployee := models.Employee{
		Id:               employeeId,
		FirstName:        "John",
		LastName:         "Doe",
		SecondLastName:   "Smith",
		DateOfBirth:      time.Date(1990, time.January, 01, 8, 0, 0, 0, time.UTC),
		DateOfEmployment: time.Date(2000, time.January, 01, 8, 0, 0, 0, time.UTC),
		Status:           "Active",
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(kafkaMessage, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()
	_ = updatedEmployee
	mockEmployeeService.On("UpdateEmployee", employeeId, employeeRequest).Return(&models.Employee{}, nil)

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeUpdate()

	assert.Error(t, err) // Will error on second ReadMessage to stop the loop
	mockConsumer.AssertCalled(t, "SubscribeTopics", mock.Anything, mock.Anything)
	mockEmployeeService.AssertCalled(t, "UpdateEmployee", employeeId, employeeRequest)
}

func TestKafkaConsumerService_ListenEmployeeUpdate_ContinuesOnInvalidJSON(t *testing.T) {
	t.Skip("skiped until mockConsumer matches with expected type")
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaMessage := &kafka.Message{
		Value: []byte("invalid json"),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(kafkaMessage, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeUpdate()

	assert.Error(t, err) // Will error on second ReadMessage to stop the loop
	mockConsumer.AssertCalled(t, "SubscribeTopics", mock.Anything, mock.Anything)
	mockEmployeeService.AssertNotCalled(t, "UpdateEmployee")
}

func TestKafkaConsumerService_ListenEmployeeUpdate_ContinuesOnServiceError(t *testing.T) {
	t.Skip("skiped until mockConsumer matches with expected type")
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	employeeId := "123"
	employeeRequest := dto.EmployeeRequest{
		FirstName:        "John",
		LastName:         "Doe",
		SecondLastName:   "Smith",
		DateOfBirth:      time.Date(1990, time.January, 01, 8, 0, 0, 0, time.UTC),
		DateOfEmployment: time.Date(2000, time.January, 01, 8, 0, 0, 0, time.UTC),
		Status:           "Active",
	}

	messageBytes := createEmployeeMessageBytes(t, employeeId, employeeRequest)
	kafkaMessage := &kafka.Message{
		Value: messageBytes,
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(kafkaMessage, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()
	mockEmployeeService.On("UpdateEmployee", employeeId, employeeRequest).Return(models.Employee{}, errors.New("service error"))

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeUpdate()

	assert.Error(t, err) // Will error on second ReadMessage to stop the loop
	mockEmployeeService.AssertCalled(t, "UpdateEmployee", employeeId, employeeRequest)
}

func TestKafkaConsumerService_ListenEmployeeDeletion_NoErrorsWhenConsumerIsDisabled(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.NoError(t, err)
	mockConsumer.AssertNotCalled(t, "SubscribeTopics")
	mockEmployeeService.AssertNotCalled(t, "DeleteEmployeeById")
}

func TestKafkaConsumerService_ListenEmployeeDeletion_ReturnsErrorWhenReadMessageFails(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	expectedError := errors.New("kafka read error")
	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, expectedError)

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockConsumer.AssertCalled(t, "SubscribeTopics", mock.Anything, mock.Anything)
	mockEmployeeService.AssertNotCalled(t, "DeleteEmployeeById")
}

func TestKafkaConsumerService_ListenEmployeeDeletion_SuccessfullyDeletesEmployee(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	employeeId := "123"
	kafkaMessage := &kafka.Message{
		Value: []byte(employeeId),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(kafkaMessage, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()
	mockEmployeeService.On("DeleteEmployeeById", employeeId).Return(nil)

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.Error(t, err) // Will error on second ReadMessage to stop the loop
	mockConsumer.AssertCalled(t, "SubscribeTopics", mock.Anything, mock.Anything)
	mockEmployeeService.AssertCalled(t, "DeleteEmployeeById", employeeId)
}

func TestKafkaConsumerService_ListenEmployeeDeletion_ContinuesOnServiceError(t *testing.T) {
	t.Skip("skiped until mockConsumer matches with expected type")
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "true")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	employeeId := "123"
	kafkaMessage := &kafka.Message{
		Value: []byte(employeeId),
	}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).Return(nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(kafkaMessage, nil).Once()
	mockConsumer.On("ReadMessage", time.Duration(-1)).Return(nil, errors.New("stop loop")).Once()
	mockEmployeeService.On("DeleteEmployeeById", employeeId).Return(errors.New("service error"))

	kafkaConsumerService := NewKafkaConsumerService(
		mockConsumer,
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.Error(t, err) // Will error on second ReadMessage to stop the loop
	mockEmployeeService.AssertCalled(t, "DeleteEmployeeById", employeeId)
}
