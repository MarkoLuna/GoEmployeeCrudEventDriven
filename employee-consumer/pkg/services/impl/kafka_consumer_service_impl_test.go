package impl

import (
	"os"
	"testing"

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

func (m *MockKafkaConsumer) ReadMessage(timeout int) (*kafka.Message, error) {
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

func TestKafkaConsumerService_ListenEmployeeInsert_NoErrorsWhenConsumerIsDisabled(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaConsumerService := NewKafkaConsumerService(
		&kafka.Consumer{},
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeInsert()

	assert.NoError(t, err)
	mockConsumer.AssertNotCalled(t, "SubscribeTopics")
	mockEmployeeService.AssertNotCalled(t, "CreateEmployee")
}

func TestKafkaConsumerService_ListenEmployeeUpdate_NoErrorsWhenConsumerIsDisabled(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaConsumerService := NewKafkaConsumerService(
		&kafka.Consumer{},
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeUpdate()

	assert.NoError(t, err)
	mockConsumer.AssertNotCalled(t, "SubscribeTopics")
	mockEmployeeService.AssertNotCalled(t, "UpdateEmployee")
}

func TestKafkaConsumerService_ListenEmployeeDeletion_NoErrorsWhenConsumerIsDisabled(t *testing.T) {
	mockConsumer := new(MockKafkaConsumer)
	mockEmployeeService := new(MockEmployeeService)

	os.Setenv("KAFKA_CONSUMER_ENABLED", "false")
	defer os.Unsetenv("KAFKA_CONSUMER_ENABLED")

	kafkaConsumerService := NewKafkaConsumerService(
		&kafka.Consumer{},
		mockEmployeeService,
	)

	err := kafkaConsumerService.ListenEmployeeDeletion()

	assert.NoError(t, err)
	mockConsumer.AssertNotCalled(t, "SubscribeTopics")
	mockEmployeeService.AssertNotCalled(t, "DeleteEmployeeById")
}
