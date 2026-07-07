package clients

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateBackoff_Exponential_Attempt1(t *testing.T) {
	result := calculateBackoff(1, BackoffStrategyExponential)
	expectedMin := 100 * time.Millisecond
	expectedMax := 100*time.Millisecond + 50*time.Millisecond
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestCalculateBackoff_Exponential_Attempt2(t *testing.T) {
	result := calculateBackoff(2, BackoffStrategyExponential)
	expectedMin := 200 * time.Millisecond
	expectedMax := 200*time.Millisecond + 100*time.Millisecond
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestCalculateBackoff_Exponential_Attempt3(t *testing.T) {
	result := calculateBackoff(3, BackoffStrategyExponential)
	expectedMin := 400 * time.Millisecond
	expectedMax := 400*time.Millisecond + 200*time.Millisecond
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestCalculateBackoff_Exponential_MaxBackoff(t *testing.T) {
	result := calculateBackoff(10, BackoffStrategyExponential)
	expectedMin := maxBackoff
	expectedMax := maxBackoff + maxBackoff/2
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestCalculateBackoff_Linear_Attempt1(t *testing.T) {
	result := calculateBackoff(1, BackoffStrategyLinear)
	expectedMin := 100 * time.Millisecond
	expectedMax := 100*time.Millisecond + 50*time.Millisecond
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestCalculateBackoff_Linear_Attempt2(t *testing.T) {
	result := calculateBackoff(2, BackoffStrategyLinear)
	expectedMin := 200 * time.Millisecond
	expectedMax := 200*time.Millisecond + 100*time.Millisecond
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestCalculateBackoff_Linear_MaxBackoff(t *testing.T) {
	result := calculateBackoff(25, BackoffStrategyLinear)
	expectedMin := maxBackoff
	expectedMax := maxBackoff + maxBackoff/2
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestCalculateBackoff_DefaultStrategy(t *testing.T) {
	result := calculateBackoff(1, BackoffStrategy(99))
	expectedMin := initialBackoff
	expectedMax := initialBackoff + initialBackoff/2
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestJitter_ZeroOrPositive(t *testing.T) {
	for i := 0; i < 50; i++ {
		j := jitter(200 * time.Millisecond)
		assert.GreaterOrEqual(t, j, time.Duration(0))
	}
}

func TestJitter_MaxValue(t *testing.T) {
	for i := 0; i < 50; i++ {
		j := jitter(200 * time.Millisecond)
		assert.Less(t, j, 100*time.Millisecond)
	}
}

func body(s string) io.ReadCloser {
	return io.NopCloser(bytes.NewReader([]byte(s)))
}

func TestDoWithRetry_SuccessFirstAttempt(t *testing.T) {
	var callCount int
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		callCount++
		return &http.Response{StatusCode: http.StatusOK, Body: body("ok")}, nil
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := doWithRetry(context.Background(), fn, req, BackoffStrategyExponential)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, callCount)
}

func TestDoWithRetry_RetryOnServerError(t *testing.T) {
	var callCount int
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			return &http.Response{StatusCode: http.StatusInternalServerError, Body: body("err")}, nil
		}
		return &http.Response{StatusCode: http.StatusOK, Body: body("ok")}, nil
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := doWithRetry(context.Background(), fn, req, BackoffStrategyExponential)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, callCount)
}

func TestDoWithRetry_AllRetriesExhausted(t *testing.T) {
	var callCount int
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		callCount++
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: body("err")}, nil
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := doWithRetry(context.Background(), fn, req, BackoffStrategyExponential)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, maxRetries+1, callCount)
}

func TestDoWithRetry_ContextCancelled(t *testing.T) {
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: body("err")}, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, _ := http.NewRequest("GET", "/test", nil)
	_, err := doWithRetry(ctx, fn, req, BackoffStrategyExponential)

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestDoWithRetry_ContextCancelledDuringBackoff(t *testing.T) {
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: body("err")}, nil
	})

	ctx, cancel := context.WithCancel(context.Background())

	req, _ := http.NewRequest("GET", "/test", nil)
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := doWithRetry(ctx, fn, req, BackoffStrategyExponential)
	assert.ErrorIs(t, err, context.Canceled)
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestDoWithRetry_RetryOnTimeout(t *testing.T) {
	var callCount int
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			return nil, timeoutError{}
		}
		return &http.Response{StatusCode: http.StatusOK, Body: body("ok")}, nil
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := doWithRetry(context.Background(), fn, req, BackoffStrategyExponential)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, callCount)
}

func TestDoWithRetry_NonRetryableErrorOnLastAttempt(t *testing.T) {
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	_, err := doWithRetry(context.Background(), fn, req, BackoffStrategyExponential)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

func TestDoWithRetry_NetOpErrorTriggersRetry(t *testing.T) {
	var callCount int
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			return nil, &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}
		}
		return &http.Response{StatusCode: http.StatusOK, Body: body("ok")}, nil
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	resp, err := doWithRetry(context.Background(), fn, req, BackoffStrategyExponential)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, callCount)
}

func TestDoWithRetry_BodyReReadable(t *testing.T) {
	bodyContent := []byte(`{"key": "value"}`)
	var callCount int
	fn := httpDoer(func(req *http.Request) (*http.Response, error) {
		callCount++
		if callCount > 1 {
			readBody, _ := io.ReadAll(req.Body)
			req.Body.Close()
			assert.Equal(t, bodyContent, readBody, "body should be re-readable on retry")
		}

		if callCount == 1 {
			return &http.Response{StatusCode: http.StatusInternalServerError, Body: body("err")}, nil
		}
		return &http.Response{StatusCode: http.StatusCreated, Body: body("ok")}, nil
	})

	req, _ := http.NewRequest("POST", "/test", io.NopCloser(bytes.NewReader(bodyContent)))
	resp, err := doWithRetry(context.Background(), fn, req, BackoffStrategyExponential)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, 2, callCount)
}
