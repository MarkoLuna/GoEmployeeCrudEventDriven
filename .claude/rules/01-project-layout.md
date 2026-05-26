# Rule: Standard Go Project Layout

> **Applies to**: all Go files in `employee-service/`, `employee-consumer/`, `common/`

## Directory Responsibilities

| Directory | Purpose |
|---|---|
| `cmd/api-server/main.go` | Entry point **only** — construct `Application{}`, call `Init()`, call `Run()`. Zero business logic. |
| `internal/app/` | Application struct, manual DI wiring, server lifecycle. |
| `internal/config/` | One file per infrastructure concern (`db_connection.go`, `kafka_configuration.go`, `auth_config.go`, `cors_config.go`, `logging_config.go`). |
| `internal/controllers/` | HTTP handlers — bind, validate, delegate to service, respond. No DB or Kafka calls. |
| `internal/services/` | Service interfaces declared here; implementations in `impl/`; test stubs in `stubs/`. |
| `internal/repositories/` | Repository interfaces + implementations + stubs. |
| `internal/models/` | DB entity structs only. No business logic. |
| `internal/routes/` | Route registration helpers, nothing else. |
| `docs/` | **Auto-generated** by `swaggo/swag`. Never edit by hand. |
| `common/` | Shared DTOs, constants, auth interface, utils — imported via local `replace` directive. |
| `pkg/` | Leave empty; reserved for code safe to publish externally. |

## Hard Rules

- **Never** import `internal/` packages across service boundaries.
- **Never** put business logic in `cmd/`.
- **Never** hand-edit files under `docs/` — always regenerate with `swag init --dir cmd/api-server,internal`.
- Cross-service shared types go in `common/` only.
- After changing `common/`, run `go mod tidy` in **all** consuming services.
