package config

import (
	"net/http"
	"time"
)

func NewHttpClient() *http.Client {
	// Create a custom client with a 10-second timeout.
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return client
}
