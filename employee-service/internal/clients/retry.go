package clients

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
)

const (
	maxRetries     = 3
	initialBackoff = 100 * time.Millisecond
	maxBackoff     = 2 * time.Second
)

// BackoffStrategy defines the algorithm for calculating retry backoff intervals.
type BackoffStrategy int

const (
	// BackoffStrategyExponential doubles the interval on each retry
	// (100ms, 200ms, 400ms, ...) up to maxBackoff.
	BackoffStrategyExponential BackoffStrategy = iota
	// BackoffStrategyLinear increases the interval by a fixed step
	// on each retry (100ms, 200ms, 300ms, ...) up to maxBackoff.
	BackoffStrategyLinear
)

// calculateBackoff returns the duration to sleep before the given retry attempt
// using the specified strategy, with random jitter added.
func calculateBackoff(attempt int, strategy BackoffStrategy) time.Duration {
	var backoff time.Duration

	switch strategy {
	case BackoffStrategyExponential:
		backoff = initialBackoff * (1 << (attempt - 1))
	case BackoffStrategyLinear:
		backoff = initialBackoff * time.Duration(attempt)
	default:
		backoff = initialBackoff
	}

	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	return backoff + jitter(backoff)
}

// jitter returns a random duration in [0, backoff/2) to spread retry timing
// and avoid thundering herd.
func jitter(backoff time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(backoff / 2)))
}

// doWithRetry executes fn with retries on transient failures (5xx status codes,
// timeouts, and net.OpError). It applies the specified backoff strategy between
// retries and respects context cancellation. The request body is buffered so it
// can be re-sent on retry.
func doWithRetry(ctx context.Context, fn httpDoer, req *http.Request, strategy BackoffStrategy) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			time.Sleep(calculateBackoff(attempt, strategy))

			if bodyBytes != nil {
				req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		resp, err = fn(req)
		if err == nil {
			if resp.StatusCode < 500 {
				return resp, nil
			}
			resp.Body.Close()
		} else {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if _, ok := err.(*net.OpError); ok {
				continue
			}
			if attempt == maxRetries {
				return nil, err
			}
			continue
		}
	}

	return resp, err
}

// httpDoer is a function that issues an HTTP request and returns a response.
type httpDoer func(*http.Request) (*http.Response, error)
