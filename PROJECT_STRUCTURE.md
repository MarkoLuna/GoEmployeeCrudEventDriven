# Project Structure - GoEmployeeCrudEventDriven

This repository follows a microservices architecture using the standard Go project layout, with a centralized shared module for core logic.

## Core Structure

```bash
GoEmployeeCrudEventDriven/
├── common/                # Shared module (DTOs, Utils, Auth)
│   ├── constants/         # Shared constants (e.g., EmployeeStatus)
│   ├── dto/               # Shared Data Transfer Objects
│   ├── services/          # Shared service interfaces and implementations
│   └── utils/             # Shared utility functions
├── employee-service/      # Main CRUD microservice (Producer)
│   ├── cmd/               # Entry points (api-server)
│   ├── docs/              # Swagger documentation
│   ├── internal/          # Private service logic (Controllers, Repositories)
│   └── pkg/               # Publicly exportable packages (if any)
├── employee-consumer/     # Background worker service (Consumer)
│   ├── cmd/               # Entry points (api-server)
│   ├── docs/              # Swagger documentation
│   ├── internal/          # Private consumer logic (Kafka listeners)
│   └── pkg/               # Publicly exportable packages (if any)
├── docker/                # Dockerfiles and orchestration configs
├── k8s/                   # Kubernetes manifests
└── bruno-collection/      # API testing collection
```

## Module Definitions

### `common/`
The **Shared Source of Truth**. This is a standalone Go module containing code that must be consistent across all services. 
- **Location**: `/common`
- **Import Path**: `github.com/MarkoLuna/GoEmployeeCrudEventDriven/common`

### `employee-service/`
Responsible for the REST API and producing Kafka events when employees are created or updated.
- **Internal Pattern**: Uses `controllers` -> `services` -> `repositories`.
- **Integration**: Imports `common/dto` for request/response models.

### `employee-consumer/`
Responsible for consuming Kafka events and performing downstream processing (e.g., database synchronization, notifications).
- **Internal Pattern**: Uses `listeners` and `handlers`.
- **Integration**: Imports `common/dto` to ensure message compatibility with the producer.

## Key Principles

- **Separation of Concerns**: Each microservice manages its own domain logic and data lifecycle.
- **Shared Contracts**: The `common` module ensures that both Producers and Consumers speak the same language (DTOs).
- **Go Project Layout**: Each service follows the `cmd/`, `internal/`, `pkg/` pattern to ensure clean encapsulation and prevent leaking internal implementation details.
- **Dependency Injection**: Services and repositories are decoupled via interfaces to facilitate testing and flexibility.