# GoEmployeeCrudEventDriven — AGENTS.md

This file serves as the centralized AI context, guidelines, and system instructions for all AI assistants (Antigravity/Gemini, Windsurf, Claude) working on this repository.

---

## Project Overview

Employee CRUD REST API built with **Go** and **Event-Driven Architecture**, split into two microservices that communicate over **Apache Kafka**.

| Service | Module | Port | Role |
|---|---|---|---|---|
| `employee-service` | `github.com/MarkoLuna/EmployeeService` | `8080` | REST API — CRUD operations, Kafka **Producer** |
| `employee-consumer` | `github.com/MarkoLuna/EmployeeConsumer` | `8081` | Background worker — Kafka **Consumer** (reads proxied from service) |
| `auth-service` | `github.com/MarkoLuna/AuthService` | `8082` | Authentication, token generation, user CRUD |
| `common` | `github.com/MarkoLuna/GoEmployeeCrudEventDriven/common` | — | Shared DTOs, Auth interface, Utils |

The `common/` module is referenced via a local `replace` directive in each service's `go.mod`:
```
replace github.com/MarkoLuna/GoEmployeeCrudEventDriven/common => ../common
```

---

## Repository Layout

```
GoEmployeeCrudEventDriven/
├── common/                ← Standalone shared Go module
├── auth-service/          ← Authentication & user management microservice
├── employee-service/      ← Kafka Producer microservice
├── employee-consumer/     ← Kafka Consumer microservice
├── docker/                ← Shared infrastructure (Keycloak compose)
├── k8s/                   ← Kubernetes manifests
├── bruno-collection/      ← Bruno API test collection
└── PROJECT_STRUCTURE.md   ← Full annotated file tree
```

See [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) for the complete annotated file tree.

---

## Internal Architecture

### Layer Pattern (both services)

```
controllers/ → services/ → repositories/ → PostgreSQL
                  ↓
            services/impl/kafka_*_service  ←→  Apache Kafka
```

### Dependency Injection Pattern
- All dependencies are expressed as **interfaces** in `internal/services/` and `internal/repositories/`.
- Concrete implementations live in `internal/services/impl/`.
- Test stubs live in `internal/services/stubs/` and `internal/repositories/`.
- Constructors return **interfaces**, not concrete types, and accept interface parameters.
- DI wiring is handled in `internal/app/application.go` per service.
- **No mocking frameworks** — test stubs are hand-written. Repository tests use `go-sqlmock`.

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
| Circuit Breaker | `github.com/sony/gobreaker` |
| Testing | `github.com/stretchr/testify`, `github.com/DATA-DOG/go-sqlmock` |
| Go version | `1.25.5` |
| CGO | **Enabled** (required by `confluent-kafka-go`) |

---

## Development Commands & Tasks

### Commands (run from each service dir)

| Task | Command |
|------|---------|
| Build+run | `make run` (or `go run cmd/api-server/main.go`) |
| Unit tests | `make test` |
| Coverage | `make test-cover` |
| Go vet | `make vet` |
| Docker Compose | `make docker-compose-run` |
| Kubernetes deploy | `make k8-apply` |
| Regenerate Swagger | `swag init --dir cmd/api-server,internal` (add `--parseDependency` for auth-service to resolve `common/dto` types) |
| gRPC / codegen | None — no codegen in this repo |

*Note: After changing `common/`, run `go mod tidy` in all four modules (`common/`, `auth-service/`, `employee-service/`, `employee-consumer/`).*

### Shared Infrastructure Tasks

- **Start Keycloak (shared)**:
  ```bash
  cd docker && docker compose -f keycloak-compose.yml up -d
  ```
- **Kubernetes (local ingress at api.employee.local)**:
  ```bash
  # Add to /etc/hosts: 127.0.0.1 api.employee.local
  kubectl apply -f k8s/
  ```

---

## Coding Conventions & Style Guidelines

### General
- Follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).
- Use `cmd/api-server/main.go` as the single entry point per service.
- Keep `internal/` strictly private — never import across service boundaries.
- Import shared types only from `common/` (DTOs, auth interfaces, utils).

### Naming
- Files: `snake_case.go`
- Interfaces: declared in their own file (e.g., `employee_service.go`).
- Implementations: `*_impl.go` (e.g., `employee_service_impl.go`).
- Stubs / mocks for tests: `*_stub.go` (located in `internal/services/stubs/` and `internal/repositories/`).
- Tests: `*_test.go` co-located with the file under test.

### Error Handling
- Return/wrap errors up the call stack; handle at the controller level for consistent HTTP status mapping.
- Use `fmt.Errorf("context: %w", err)` — never `panic` or swallow errors in service/repository code.
- Use structured logging (configured in `internal/config/logging_config.go`).
- Log with struct/method context for observability (e.g., log struct and method names).

### Context Propagation
- Ensure all HTTP calls, database operations, and background worker processing respect context propagation (`context.Context`).

### Security
- For every HTTP service call, ensure the JWT token is injected using the builder pattern.
- Set `OAUTH_ENABLED=false` for local dev without Keycloak (configured in each service's `.env`).

### Testing
- Unit tests must use hand-written stubs (not real DB or Kafka connections).
- Integration/Repository tests use `go-sqlmock`.
- Run tests with `make test` from inside each service directory.
- Maintain high test coverage for business logic.

---

## Environment Variables

Each service reads `.env` from its own working directory.

| Variable | Default | Notes |
|---|---|---|
| `SERVER_PORT` | `8080` / `8081` | HTTP listen port |
| `SERVER_HOST` | `0.0.0.0` | HTTP listen host |
| `SERVER_SSL_ENABLED` | `false` | TLS toggle |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` / `5433` | PostgreSQL port (service DB=`5432`, consumer DB=`5433`, auth DB=`5434` when local without Docker) |
| `DB_NAME` | `employee_db` | Database name |
| `DB_USER` | `employee_user` | Database user |
| `DB_PASSWORD` | `employeepw` | Database password |
| `DB_DRIVER_NAME` | `postgres` | Go database driver |
| `KAFKA_BOOTSTRAP_SERVERS` | `localhost:9092` | Kafka broker address |
| `KAFKA_CONSUMER_GROUP_ID` | `employee-group` | Consumer group (consumer only) |
| `HTTP_TIMEOUT` | `30s` | HTTP client timeout for consumer proxy calls |
| `OAUTH_ENABLED` | `false` | Keycloak auth toggle |

---

## Key API Endpoints (auth-service :8082)

| Method | Path | Description |
|---|---|---|
| `POST` | `/oauth/token` | Generate JWT token |
| `GET` | `/oauth/userinfo` | Get token claims |
| `GET` | `/api/user/` | List all users |
| `GET` | `/api/user/:id` | Get user by ID |
| `POST` | `/api/user/` | Create user |
| `PUT` | `/api/user/:id` | Update user |
| `DELETE` | `/api/user/:id` | Delete user |

## Key API Endpoints (employee-service :8080)

| Method | Path | Description |
|---|---|---|
| `GET` | `/healthcheck/` | Liveness check |
| `GET` | `/swagger/` | Swagger UI |
| `GET` | `/api/employee/` | List all employees (proxied to consumer) |
| `GET` | `/api/employee/:id` | Get employee by ID (proxied to consumer) |
| `POST` | `/api/employee/` | Create employee (triggers Kafka event) |
| `PUT` | `/api/employee/:id` | Update employee (triggers Kafka event) |
| `DELETE` | `/api/employee/:id` | Delete employee (triggers Kafka event) |

---

## Quirks & Gotchas

- **CGO must be enabled** — `confluent-kafka-go` requires a C compiler (`CGO_ENABLED=1`). Docker build disables CGO (`CGO_ENABLED=0`) for static linking.
- **Never hand-edit `docs/`** — Swagger docs are code-generated from controller annotations via `swag init`.
- **Kafka consumer has DLT + exponential backoff** — dead-letter topics (`*.v1.dlt`) catch final failures. Idempotency via `sync.Map` per message key.
- **Kafka topic naming** — `<entity>-<action>.v<N>`, DLT: `<topic>.dlt`.
- **Consumer graceful shutdown** uses `signal.NotifyContext` + `WaitGroupCoordinator` — ensure new background services register with the coordinator.
- **`internal/` is strictly private** across service boundaries. Cross-service sharing goes in `common/` only.
- **`common/` local replace** — the `common/` module is a local replace, never published; edit in place.
- **Bruno collection** — `bruno-collection/` covers all REST CRUD endpoints for manual testing.
- **Use `utils.GetEnv(key, default)`** for env variables, never raw `os.Getenv`.
- **Service Communication** — GET (read) operations proxy requests from the `employee-service` to the `employee-consumer` via HTTP using `EMPLOYEE_CONSUMER_HOST`. Command/write actions propagate from the `employee-service` to the `employee-consumer` via Kafka events.
- **HTTP Circuit Breaker** — `employee-service` wraps consumer HTTP calls with sony/gobreaker (file: `internal/clients/circuit_breaker.go`). Trips at ≥60% failure rate over ≥5 requests. Prevents cascading failures.
- **HTTP Retry with Backoff** — Consumer proxy calls retry up to 3 times on 5xx/timeouts/net.OpError (file: `internal/clients/retry.go`). Supports `BackoffStrategyExponential` (default) and `BackoffStrategyLinear`. Request body is buffered for re-send on retry.
- **Kafka Producer Idempotence** — `enable.idempotence=true` in `internal/services/impl/kafka_producer_service_impl.go` ensures exactly-once publishing with ordered delivery.
- **Async 202 Pattern** — Write endpoints (POST/PUT/DELETE) return `202 Accepted` immediately after Kafka publish, decoupling from consumer processing.
