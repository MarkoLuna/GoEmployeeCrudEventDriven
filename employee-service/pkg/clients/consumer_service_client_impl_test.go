package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeService/pkg/constants"
	"github.com/MarkoLuna/EmployeeService/pkg/models"
	"github.com/stretchr/testify/assert"
)

func createTestEmployee() models.Employee {
	return models.Employee{
		Id:               "123",
		FirstName:        "John",
		LastName:         "Doe",
		SecondLastName:   "Smith",
		DateOfBirth:      time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		DateOfEmployment: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		Status:           constants.ACTIVE,
	}
}

func TestEmployeeConsumerServiceClientImpl_Create_Success(t *testing.T) {
	expectedEmployee := createTestEmployee()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/employee/", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedEmployee)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	result, err := serviceClient.Create(expectedEmployee)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmployee.Id, result.Id)
	assert.Equal(t, expectedEmployee.FirstName, result.FirstName)
	assert.Equal(t, expectedEmployee.LastName, result.LastName)
}

func TestEmployeeConsumerServiceClientImpl_Create_ErrorFromService(t *testing.T) {
	validEmployee := models.Employee{
		Id: "test",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	_, err := serviceClient.Create(validEmployee)
	assert.Error(t, err)
}

func TestEmployeeConsumerServiceClientImpl_FindAll_Success(t *testing.T) {
	employees := []models.Employee{
		createTestEmployee(),
		{
			Id:               "456",
			FirstName:        "Jane",
			LastName:         "Doe",
			SecondLastName:   "Smith",
			DateOfBirth:      time.Date(1992, time.February, 2, 0, 0, 0, 0, time.UTC),
			DateOfEmployment: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Status:           constants.ACTIVE,
		},
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/employee/", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employees)
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	result, err := serviceClient.FindAll()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, employees[0].Id, result[0].Id)
	assert.Equal(t, employees[1].Id, result[1].Id)
}

func TestEmployeeConsumerServiceClientImpl_FindAll_EmptyList(t *testing.T) {
	// Create mock server returning empty list
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.Employee{})
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	result, err := serviceClient.FindAll()

	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))
}

func TestEmployeeConsumerServiceClientImpl_FindById_Success(t *testing.T) {
	expectedEmployee := createTestEmployee()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/employee/123", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedEmployee)
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	result, err := serviceClient.FindById("123")

	assert.NoError(t, err)
	assert.Equal(t, expectedEmployee.Id, result.Id)
	assert.Equal(t, expectedEmployee.FirstName, result.FirstName)
	assert.Equal(t, expectedEmployee.LastName, result.LastName)
}

func TestEmployeeConsumerServiceClientImpl_FindById_NotFound(t *testing.T) {
	// Create mock server returning 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	_, err := serviceClient.FindById("999")

	assert.Error(t, err)
}

func TestEmployeeConsumerServiceClientImpl_DeleteById_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/api/employee/123", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	err := serviceClient.DeleteById("123")

	assert.NoError(t, err)
}

func TestEmployeeConsumerServiceClientImpl_DeleteById_NotFound(t *testing.T) {
	// Create mock server returning 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	err := serviceClient.DeleteById("999")

	assert.NoError(t, err)
}

func TestEmployeeConsumerServiceClientImpl_Update_Success(t *testing.T) {
	employeeToUpdate := createTestEmployee()
	employeeToUpdate.FirstName = "UpdatedName"

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/api/employee/123", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employeeToUpdate)
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	result, err := serviceClient.Update(employeeToUpdate)

	assert.NoError(t, err)
	assert.Equal(t, employeeToUpdate.Id, result.Id)
	assert.Equal(t, "UpdatedName", result.FirstName)
}

func TestEmployeeConsumerServiceClientImpl_Update_NotFound(t *testing.T) {
	employee := createTestEmployee()

	// Create mock server returning 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := http.Client{}
	serviceClient := &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: server.URL,
	}

	_, err := serviceClient.Update(employee)

	assert.Error(t, err)
}

func TestNewEmployeeConsumerServiceClient(t *testing.T) {
	// Test the constructor
	client := http.Client{}
	serviceClient := NewEmployeeConsumerServiceClient(client)

	assert.NotNil(t, serviceClient)

	impl, ok := serviceClient.(*EmployeeConsumerServiceClientImpl)
	assert.True(t, ok)
	assert.NotNil(t, impl)
	assert.Equal(t, employeeConsumerHost, impl.serviceHost)
}
