package clients

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

var (
	consumerCircuitBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "employee-consumer-http",
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("circuit breaker '%s' state changed: %s -> %s", name, from, to)
		},
	})
)

// circuitBreakerRoundTripper wraps an http.RoundTripper with gobreaker
// to prevent cascading failures when the downstream service is unhealthy.
type circuitBreakerRoundTripper struct {
	next http.RoundTripper
}

// RoundTrip executes the request through the circuit breaker.
// 5xx responses and connection errors are recorded as failures.
func (c *circuitBreakerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	result, err := consumerCircuitBreaker.Execute(func() (interface{}, error) {
		resp, err := c.next.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			return nil, errors.New("server error")
		}
		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*http.Response), nil
}

// NewCircuitBreakerTransport returns an http.RoundTripper wrapping
// http.DefaultTransport with a circuit breaker for employee-consumer
// HTTP calls. The breaker trips when the failure rate exceeds 60%
// over at least 5 requests in a 60-second window.
func NewCircuitBreakerTransport() http.RoundTripper {
	return &circuitBreakerRoundTripper{
		next: http.DefaultTransport,
	}
}
