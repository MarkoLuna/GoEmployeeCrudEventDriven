# GoEmployeeCrudEventDriven — Gemini Project Context

## Skills

| Skill | File | When to apply |
|---|---|---|
| Go Backend Development | [skills/golang_backend.md](skills/golang_backend.md) | Any Go code change in this repo |

---

## Project Overview

Employee CRUD REST API built with **Go** and **Event-Driven Architecture**, split into two microservices that communicate over **Apache Kafka**.

| Service | Module | Port | Role |
|---|---|---|---|
| `employee-service` | `github.com/MarkoLuna/EmployeeService` | `8080` | REST API — CRUD operations, Kafka **Producer** |
| `employee-consumer` | `github.com/MarkoLuna/EmployeeConsumer` | `8081` | Background worker — Kafka **Consumer** |
| `common` | `github.com/MarkoLuna/GoEmployeeCrudEventDriven/common` | — | Shared DTOs, Auth interface, Utils |

The `common/` module is referenced via a local `replace` directive in each service's `go.mod`:
```
replace github.com/MarkoLuna/GoEmployeeCrudEventDriven/common => ../common
```

---

## Repository Layout

```
GoEmployeeCrudEventDriven/
├── .gemini/               ← You are here
├── common/                ← Standalone shared Go module
├── employee-service/      ← Kafka Producer microservice
├── employee-consumer/     ← Kafka Consumer microservice
├── docker/                ← Shared infrastructure (Keycloak compose)
├── k8s/                   ← Kubernetes manifests
├── bruno-collection/      ← Bruno API test collection
└── PROJECT_STRUCTURE.md   ← Full annotated file tree
```

See [PROJECT_STRUCTURE.md](../PROJECT_STRUCTURE.md) for the complete annotated file tree.

---

## Internal Architecture

### Layer Pattern (both services)

```
controllers/ → services/ → repositories/ → PostgreSQL
                  ↓
            services/impl/kafka_*_service  ←→  Apache Kafka
```

### Dependency Injection Pattern

- All dependencies are expressed as **interfaces** in `internal/services/` and `internal/repositories/`
- Concrete implementations live in `internal/services/impl/`
- Test stubs live in `internal/services/stubs/` and `internal/repositories/`
- The application is wired in `internal/app/application.go`

---

## Technology Stack

| Concern | Library |
|---|---|
| HTTP Framework | `github.com/labstack/echo/v4` |
| Kafka | `github.com/confluentinc/confluent-kafka-go` (CGO required) |
| PostgreSQL | `github.com/lib/pq` |
| JWT / OAuth | `github.com/golang-jwt/jwt`, `github.com/go-oauth2/oauth2/v4` |
| Keycloak | Custom wrapper in `common/services/auth/keycloak_impl.go` |
| Validation | `gopkg.in/go-playground/validator.v9` |
| Swagger | `github.com/swaggo/swag` + `github.com/swaggo/echo-swagger` |
| Testing | `github.com/stretchr/testify`, `github.com/DATA-DOG/go-sqlmock` |
| Go version | `1.25.5` |
| CGO | **Enabled** (required by `confluent-kafka-go`) |

---

## Coding Conventions

### General
- Follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- Use `cmd/api-server/main.go` as the single entry point per service
- Keep `internal/` strictly private — never import across services
- Import shared types only from `common/` (DTOs, auth interfaces, utils)

### Naming
- Files: `snake_case.go`
- Interfaces: declared in their own file (e.g., `employee_service.go`)
- Implementations: `*_impl.go` (e.g., `employee_service_impl.go`)
- Stubs / mocks for tests: `*_stub.go`
- Tests: `*_test.go` co-located with the file under test

### Error Handling
- Return errors up the call stack; handle at the controller level
- Use structured logging (configured in `internal/config/logging_config.go`)
- Log with struct/method context for observability

### Testing
- Unit tests use stubs (not real DB/Kafka connections)
- Integration tests may use `go-sqlmock` for repository tests
- Run tests with `make test` from inside each service directory

---

## Environment Variables

Each service reads `.env` from its own working directory.

| Variable | Default | Notes |
|---|---|---|
| `SERVER_PORT` | `8080` / `8081` | HTTP listen port |
| `SERVER_HOST` | `0.0.0.0` | HTTP listen host |
| `SERVER_SSL_ENABLED` | `false` | TLS toggle |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` / `5433` | PostgreSQL port |
| `DB_NAME` | `employee_db` | Database name |
| `DB_USER` | `employee_user` | Database user |
| `DB_PASSWORD` | `employeepw` | Database password |
| `DB_DRIVER_NAME` | `postgres` | Go database driver |
| `KAFKA_BOOTSTRAP_SERVERS` | `localhost:9092` | Kafka broker address |
| `KAFKA_CONSUMER_GROUP_ID` | `employee-group` | Consumer group (consumer only) |
| `OAUTH_ENABLED` | `false` | Keycloak auth toggle |

---

## Common Tasks

### Run a service locally
```bash
cd employee-service   # or employee-consumer
make run
# or directly:
go run cmd/api-server/main.go
```

### Run tests
```bash
cd employee-service && make test
cd employee-consumer && make test
```

### Regenerate Swagger docs
```bash
# Run from inside the service directory
swag init --dir cmd/api-server,internal
```

### Run with Docker Compose
```bash
cd employee-service && make docker-compose-run
cd employee-consumer && make docker-compose-run
```

### Start Keycloak (shared)
```bash
cd docker && docker compose -f keycloak-compose.yml up -d
```

### Kubernetes (local ingress at api.employee.local)
```bash
# Add to /etc/hosts: 127.0.0.1 api.employee.local
kubectl apply -f k8s/
```

---

## Key API Endpoints (employee-service :8080)

| Method | Path | Description |
|---|---|---|
| `GET` | `/healthcheck/` | Liveness check |
| `GET` | `/swagger/` | Swagger UI |
| `GET` | `/api/employee/` | List all employees |
| `GET` | `/api/employee/:id` | Get employee by ID |
| `POST` | `/api/employee/` | Create employee (triggers Kafka event) |
| `PUT` | `/api/employee/:id` | Update employee (triggers Kafka event) |
| `DELETE` | `/api/employee/:id` | Delete employee |

---

## Important Notes

- **CGO must be enabled** — `confluent-kafka-go` requires a C compiler (GCC/Clang).
- The `common/` module is a **local replace** — it is never published to a registry; always edit in place.
- Swagger annotations live in controller files and are code-generated into `docs/`.
- Bruno collection (`bruno-collection/`) covers all CRUD endpoints and can be used for manual API testing.
