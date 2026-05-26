# Rule: Kafka Patterns

> **Applies to**: `employee-service/` (Producer) and `employee-consumer/` (Consumer).

## Interface Wrappers for Kafka Clients

Always wrap the concrete confluent-kafka-go type behind a local interface — never depend on the struct directly:

```go
type KafkaConsumer interface {
    SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
    ReadMessage(timeout time.Duration) (*kafka.Message, error)
    Close() error
}

type KafkaProducer interface {
    Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
    Close()
}
```

This keeps services fully testable without a live broker.

## Producer: Topic Names from Environment

```go
var (
    employeeInsertTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_INSERT_TOPIC", "employee-insert.v1")
    employeeUpdateTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_UPDATE_TOPIC", "employee-update.v1")
    employeeDeleteTopic = utils.GetEnv("KAFKA_PRODUCER_EMPLOYEE_DELETE_TOPIC", "employee-deletion.v1")
)
```

Topic naming convention: `<entity>-<action>.v<N>` (e.g., `employee-insert.v1`).

## Producer: Serialize to JSON Before Producing

```go
func (kSrv KafkaProducerService) SendInsert(employee dto.EmployeeMessage) error {
    value, err := json.Marshal(employee)
    if err != nil {
        return fmt.Errorf("marshal employee message: %w", err)
    }
    return kSrv.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &employeeInsertTopic,
            Partition: kafka.PartitionAny,
        },
        Value: value,
    }, nil)
}
```

## Consumer: Fan-out Worker Pools

Route each topic to a dedicated buffered channel + fixed-size goroutine pool:

```go
insertCh := make(chan *kafka.Message, insertWorkerCount)
updateCh := make(chan *kafka.Message, updateWorkerCount)
deleteCh := make(chan *kafka.Message, deleteWorkerCount)

startWorkerPool(&wg, employeeInsertTopic, insertWorkerCount, insertCh, kSrv.handleInsert)
startWorkerPool(&wg, employeeUpdateTopic, updateWorkerCount, updateCh, kSrv.handleUpdate)
startWorkerPool(&wg, employeeDeleteTopic, deleteWorkerCount, deleteCh, kSrv.handleDelete)
```

**Never** spawn unbounded goroutines. Always use a fixed-size worker pool.

## Consumer: Idempotency with `sync.Map`

```go
func messageKey(msg *kafka.Message) string {
    return fmt.Sprintf("%s:%d:%d", *msg.TopicPartition.Topic,
        msg.TopicPartition.Partition, msg.TopicPartition.Offset)
}

if _, loaded := kSrv.processedKeys.LoadOrStore(key, struct{}{}); loaded {
    log.Printf("duplicate message %s, skipping", key)
    return
}
```

## Consumer: Exponential Backoff Retry

```go
func (kSrv *KafkaConsumerServiceImpl) withRetry(fn func() error) error {
    backoff := kSrv.retryInitialBackoff
    for i := 0; i < kSrv.retryMaxAttempts; i++ {
        if err := fn(); err == nil {
            return nil
        } else if i == kSrv.retryMaxAttempts-1 {
            return err
        }
        time.Sleep(backoff)
        backoff = min(backoff*2, 30*time.Second)
    }
    return nil
}
```

## Consumer: Dead Letter Topic (DLT) on Final Failure

```go
func (kSrv *KafkaConsumerServiceImpl) produceToDLT(dltTopic string, originalMsg *kafka.Message, finalErr error) {
    dltMsg := &kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &dltTopic, Partition: kafka.PartitionAny},
        Value:          originalMsg.Value,
        Key:            originalMsg.Key,
        Headers: append(originalMsg.Headers,
            kafka.Header{Key: "x-retries",        Value: []byte(strconv.Itoa(kSrv.retryMaxAttempts))},
            kafka.Header{Key: "x-error-message",  Value: []byte(finalErr.Error())},
            kafka.Header{Key: "x-original-topic", Value: []byte(*originalMsg.TopicPartition.Topic)},
        ),
    }
    kSrv.producer.Produce(dltMsg, nil)
}
```

DLT topic naming: `<original-topic>.dlt` (e.g., `employee-insert.v1.dlt`).

## CGO Requirement

`confluent-kafka-go` requires CGO. Ensure `CGO_ENABLED=1` and a C compiler (GCC/Clang) are present in any build environment.
