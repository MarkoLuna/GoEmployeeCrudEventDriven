package clients

import (
	"net/http"
)

// EmployeeConsumerServiceClientBuilder constructs an EmployeeConsumerServiceClientImpl
// using the builder pattern, allowing optional configuration such as a JWT token
// before the client instance is created.
type EmployeeConsumerServiceClientBuilder struct {
	client         http.Client
	serviceHost    string
	token          string
	customInstance EmployeeConsumerServiceClient
}

// NewEmployeeConsumerServiceClientBuilder creates a new builder with sensible defaults.
// The serviceHost is pre-populated from the package-level environment variable.
func NewEmployeeConsumerServiceClientBuilder() *EmployeeConsumerServiceClientBuilder {
	return &EmployeeConsumerServiceClientBuilder{
		client:      http.Client{},
		serviceHost: employeeConsumerHost,
		token:       "",
	}
}

// WithHttpClient sets a custom http.Client on the builder.
func (b *EmployeeConsumerServiceClientBuilder) WithHttpClient(client http.Client) *EmployeeConsumerServiceClientBuilder {
	b.client = client
	return b
}

// WithServiceHost overrides the target service host URL.
func (b *EmployeeConsumerServiceClientBuilder) WithServiceHost(host string) *EmployeeConsumerServiceClientBuilder {
	b.serviceHost = host
	return b
}

// WithJwtToken sets the initial JWT bearer token that will be injected into
// every outgoing request's Authorization header.
func (b *EmployeeConsumerServiceClientBuilder) WithJwtToken(token string) *EmployeeConsumerServiceClientBuilder {
	b.token = token
	return b
}

// WithCustomInstance allows injecting a pre-built instance (e.g. a mock or stub)
// that Build() will return. Useful for testing.
func (b *EmployeeConsumerServiceClientBuilder) WithCustomInstance(instance EmployeeConsumerServiceClient) *EmployeeConsumerServiceClientBuilder {
	b.customInstance = instance
	return b
}

// Build constructs and returns a fully configured EmployeeConsumerServiceClient.
func (b *EmployeeConsumerServiceClientBuilder) Build() EmployeeConsumerServiceClient {
	if b.customInstance != nil {
		return b.customInstance
	}
	return &EmployeeConsumerServiceClientImpl{
		client:      b.client,
		serviceHost: b.serviceHost,
		token:       b.token,
	}
}
