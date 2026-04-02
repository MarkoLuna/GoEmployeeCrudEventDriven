# Common Module - GoEmployeeCrudEventDriven

This module repository houses the shared logic, constants, and data transfer objects (DTOs) used across the microservices in the `GoEmployeeCrudEventDriven` project.

## Purpose

The `common` module is designed to eliminate code duplication and ensure consistency across the services (`employee-service` and `employee-consumer`). It provides standardized implementations for:
- **Authentication**: Unified OAuth interfaces and Keycloak RS256 validation.
- **Environment Management**: Consistent configuration parsing.
- **Data Schemas**: Shared DTOs for Kafka messages and REST requests.
- **Validation**: Shared validation rules for enterprise-level consistency.

## Usage in Services

To use this module in a new or existing service, add the following to your `go.mod`:

```go
require github.com/MarkoLuna/GoEmployeeCrudEventDriven/common v0.0.0

replace github.com/MarkoLuna/GoEmployeeCrudEventDriven/common => ../common
```

Then, import the required packages in your Go files:

```go
import "github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
import "github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
```
