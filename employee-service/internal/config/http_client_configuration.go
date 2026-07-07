package config

import (
	"net/http"
	"time"

	"github.com/MarkoLuna/EmployeeService/internal/clients"
)

func NewHttpClient() *http.Client {
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: clients.NewCircuitBreakerTransport(),
	}

	return client
}
