package services

import "context"

type KafkaConsumerService interface {
	// Listen polls Kafka and dispatches messages to per-topic worker pools.
	// It returns when ctx is cancelled (graceful shutdown) or an unrecoverable
	// broker error occurs.
	Listen(ctx context.Context) error

	// Stop signals the consumer to stop, closes the underlying Kafka consumer,
	// and waits for all in-flight worker goroutines to finish.
	Stop()
}
