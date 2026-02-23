# GoEmployeeCrudEventDriven

Employee CRUD REST API using **Go** and **Event-Driven Architecture**.  
The project is split into two services that communicate over **Apache Kafka**:

| Service | Port | Description |
|---|---|---|
| `employee-service` | `8080` | REST API — CRUD operations, publishes Kafka events |
| `employee-consumer` | `8081` | Kafka consumer — processes events, stores/syncs data |

---

## Prerequisites

### 1. Go 1.21+
```bash
go version
```

### 2. GCC / CGO (required by confluent-kafka-go)
CGO must be enabled and a C compiler must be present.

```bash
# Enable CGO
go env -w CGO_ENABLED="1"

# Verify
go env CGO_ENABLED
```

**Install GCC:**
```bash
# Ubuntu / Debian
sudo apt-get install build-essential

# macOS
xcode-select --install
```

### 3. PostgreSQL
Each service connects to a PostgreSQL database. Default connection values:

| Variable | Default |
|---|---|
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` (service) / `5433` (consumer) |
| `DB_NAME` | `employee_db` |
| `DB_USER` | `employee_user` |
| `DB_PASSWORD` | `employeepw` |
| `DB_DRIVER_NAME` | `postgres` |

The database schema is initialized automatically when using Docker Compose via `resources/init.sql`:

```sql
CREATE TABLE employees (
  id_employee       TEXT PRIMARY KEY,
  first_name        TEXT,
  last_name         TEXT,
  second_last_name  TEXT,
  date_of_birth     DATE,
  date_of_employment DATE,
  status            TEXT
);
```

### 4. Apache Kafka
Both services require a running Kafka broker.

| Variable | Default |
|---|---|
| `KAFKA_BOOTSTRAP_SERVERS` | `localhost:9092` |
| `KAFKA_CONSUMER_GROUP_ID` | `employee-group` (consumer only) |

### 5. Environment Variables (optional `.env` file)
Each service reads a `.env` file from its working directory if present.

```bash
# employee-service / employee-consumer
DB_HOST=localhost
DB_PORT=5432
DB_NAME=employee_db
DB_USER=employee_user
DB_PASSWORD=employeepw
DB_DRIVER_NAME=postgres
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
OAUTH_ENABLED=false
SERVER_SSL_ENABLED=false
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_CONSUMER_GROUP_ID=employee-group
```

---

## Running Locally

### Install dependencies (run inside each service directory)

```bash
cd employee-service && go mod tidy
cd employee-consumer && go mod tidy
```

### Run employee-service

```bash
cd employee-service

# directly
go run cmd/api-server/main.go

# or with make
make run
```

### Run employee-consumer

```bash
cd employee-consumer

# directly
go run cmd/api-server/main.go

# or with make
make run
```

---

## Running with Docker Compose

Each service has its own `docker-compose.yml` that starts the service together with PostgreSQL.

```bash
# employee-service
cd employee-service && make docker-compose-run

# employee-consumer
cd employee-consumer && make docker-compose-run
```

---

## Running Tests

```bash
# employee-service
cd employee-service && make test

# employee-consumer
cd employee-consumer && make test
```

---

## Healthcheck

```bash
# employee-service (port 8080)
curl -X GET http://localhost:8080/healthcheck/

# employee-consumer (port 8081)
curl -X GET http://localhost:8081/healthcheck/
```

---

## Swagger UI

```
http://localhost:8080/swagger/   # employee-service
http://localhost:8081/swagger/   # employee-consumer
```

To regenerate Swagger docs after changing annotations:
```bash
# Run from inside the service directory
swag init --dir cmd/api-server,internal
```

---

## Example Curls (employee-service)

```bash
# Get all employees
curl http://localhost:8080/api/employee/

# Get employee by ID
curl http://localhost:8080/api/employee/{id}

# Create employee
curl --request POST 'http://localhost:8080/api/employee/' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "firstName": "Marcos",
    "lastName": "Luna",
    "secondLastName": "Valdez",
    "dateOfBirth": "1994-04-25T12:00:00Z",
    "dateOfEmployment": "1994-04-25T12:00:00Z",
    "status": "ACTIVE"
  }'

# Update employee
curl --request PUT 'http://localhost:8080/api/employee/{id}' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "firstName": "Gerardo",
    "lastName": "Luna",
    "secondLastName": "Valdezz",
    "dateOfBirth": "1994-04-25T12:00:00Z",
    "dateOfEmployment": "0001-01-01T00:00:00Z",
    "status": "INACTIVE"
  }'

# Delete employee
curl --request DELETE 'http://localhost:8080/api/employee/{id}'
```

---

## Project Structure

Both services follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
<service>/
├── cmd/api-server/   # Entry point (main.go)
├── internal/         # Private application code
│   ├── app/
│   ├── config/
│   ├── controllers/
│   ├── services/
│   ├── repositories/
│   ├── models/
│   ├── routes/
│   ├── dto/
│   └── constants/
├── pkg/utils/        # Public utility packages
├── docs/             # Swagger generated docs
├── resources/        # SQL init scripts, SSL certs
└── Makefile
```

See [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for more details.
