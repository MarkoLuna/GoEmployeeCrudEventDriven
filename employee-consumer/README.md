# Employee Crud

Employee Crud Rest and Kafka Consumer API using Golang

## Run with docker

```bash
make docker-run
```

## Run with docker compose

```bash
make docker-compose-run
```

## Run in k8 locally

```bash
$ make k8-apply
$ kubectl exec -it employeecrud-consumer-pod -- apk add curl
$ kubectl exec -it employeecrud-consumer-pod -- curl -X GET http://employeecrud-consumer-pod:8081/healthcheck/
$ make k8-remove 
```

## Resiliency and Error Handling

The consumer service is designed for high availability and reliable message processing through the following features:

### 1. Per-Topic Concurrency
Each Kafka topic (insert, update, delete) is processed by an independent pool of workers, allowing different operation types to scale without shared bottlenecks.

### 2. Exponential Backoff Retry Policy
If a transient failure occurs (e.g., database downtime), the consumer will automatically retry the operation with an exponential backoff.
- **Configurable Attempts**: Default is 3 attempts.
- **Backoff Strategy**: Initial delay (default 500ms) that doubles on each attempt, capped at 30 seconds.
- **Safety**: Permanent errors (like malformed JSON) are skipped immediately to prevent infinite retry loops.

### 3. Dead Letter Topics (DLT)
Messages that fail all retry attempts are automatically moved to a **Dead Letter Topic**. This ensures no data is lost and allows for manual inspection of failing messages.
- **DLT Topics**: Each primary topic has its own DLT (e.g., `employee-insert.v1.dlt`).
- **Metadata**: DLT messages include debugging headers:
    - `x-retries`: Total attempts made before failure.
    - `x-error-message`: The reason for the final failure.
    - `x-original-topic`: The origin of the message.

---

## Configuration (Environment Variables)

| Variable | Description | Default |
| :--- | :--- | :--- |
| `KAFKA_BOOTSTRAP_SERVERS` | List of Kafka brokers | `localhost:9092` |
| `KAFKA_CONSUMER_GROUP_ID` | Kafka consumer group identifier | `employee-group` |
| `KAFKA_CONSUMER_ENABLED` | Toggle message consumption | `true` |
| **Concurrency** | | |
| `KAFKA_CONSUMER_MAX_WORKERS_INSERT` | Workers for insert topic | `3` |
| `KAFKA_CONSUMER_MAX_WORKERS_UPDATE` | Workers for update topic | `3` |
| `KAFKA_CONSUMER_MAX_WORKERS_DELETE` | Workers for delete topic | `3` |
| **Retries** | | |
| `KAFKA_CONSUMER_MAX_RETRIES` | Max attempts per message | `3` |
| `KAFKA_CONSUMER_RETRY_INITIAL_BACKOFF_MS` | Starting retry delay | `500` |
| **Dead Letter Topics** | | |
| `KAFKA_CONSUMER_EMPLOYEE_INSERT_DLT` | DLT for insert failures | `employee-insert.v1.dlt` |
| `KAFKA_CONSUMER_EMPLOYEE_UPDATE_DLT` | DLT for update failures | `employee-update.v1.dlt` |
| `KAFKA_CONSUMER_EMPLOYEE_DELETE_DLT` | DLT for delete failures | `employee-deletion.v1.dlt` |

---

## Example Curls

```bash
# get all employees
$ curl http://localhost:8081/api/employee/
[]

# get employee by id
$ curl http://localhost:8081/api/employee/1 

# create employee
$ curl --location --request POST 'http://localhost:8081/api/employee/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "firstName": "Marcos",
    "lastName": "Luna",
    "secondLastName": "Valdez",
    "dateOfBirth": "1994-04-25T12:00:00Z",
    "dateOfEmployment": "1994-04-25T12:00:00Z",
    "status": "ACTIVE"
}'

# delete employee
$ curl -X DELETE 'http://localhost:8081/api/employee/2'

# update employee
$ curl --location --request PUT 'http://localhost:8081/api/employee/3' \
--header 'Content-Type: application/json' \
--data-raw '{
    "firstName": "Gerardo",
    "lastName": "Luna",
    "secondLastName": "Valdezz",
    "dateOfBirth": "1994-04-25T12:00:00Z",
    "dateOfEmployment": "0001-01-01T00:00:00Z",
    "status": "INACTIVE"
}'
```
