package clients

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	fn func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fn(req)
}

func TestNewCircuitBreakerTransport_ReturnsTransport(t *testing.T) {
	transport := NewCircuitBreakerTransport()
	assert.NotNil(t, transport)
}

func TestCircuitBreakerRoundTrip_Success(t *testing.T) {
	cb := &circuitBreakerRoundTripper{
		next: &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("ok")),
				}, nil
			},
		},
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := cb.RoundTrip(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCircuitBreakerRoundTrip_ServerError(t *testing.T) {
	cb := &circuitBreakerRoundTripper{
		next: &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("error")),
				}, nil
			},
		},
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := cb.RoundTrip(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCircuitBreakerRoundTrip_NextReturnsError(t *testing.T) {
	cb := &circuitBreakerRoundTripper{
		next: &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("connection reset")
			},
		},
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := cb.RoundTrip(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "connection reset")
}

func TestCircuitBreakerRoundTrip_ClientErrorPassesThrough(t *testing.T) {
	cb := &circuitBreakerRoundTripper{
		next: &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("not found")),
				}, nil
			},
		},
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := cb.RoundTrip(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
