# Rule: Naming Conventions

> **Applies to**: all Go source files across all modules.

## File Names

| Pattern | Convention | Example |
|---|---|---|
| All Go source files | `snake_case.go` | `employee_repository_impl.go` |
| Interfaces | `<entity>_<role>.go` | `employee_service.go` |
| Implementations | `<entity>_<role>_impl.go` | `employee_service_impl.go` |
| Test stubs | `<entity>_<role>_stub.go` | `employee_service_stub.go` |
| Tests | `<subject>_test.go` co-located | `employee_service_impl_test.go` |

## Go Identifiers

| Element | Convention | Example |
|---|---|---|
| Packages | lowercase, single word | `package repositories` |
| Interfaces | PascalCase, **no** `I` prefix | `EmployeeRepository` |
| Concrete types | PascalCase + `Impl` suffix | `EmployeeRepositoryImpl` |
| Test stubs | PascalCase + `Stub` suffix | `EmployeeRepositoryStub` |
| Constructors | `New<Type>(...)` accepting interfaces | `NewEmployeeRepository(db *sql.DB, initDb bool) EmployeeRepository` |
| Value receiver names | short abbreviation of type | `(er EmployeeRepositoryImpl)` |
| Unexported helpers | camelCase | `messageKey()`, `isConsumerEnabled()` |
| Env-var defaults / topic strings | kebab-case strings | `"employee-insert.v1"` |

## Constructor Rule

Constructors must:
1. Accept **interface** parameters, not concrete types.
2. Return the **interface** type, not the concrete type.

```go
// ✅ correct
func NewEmployeeRepository(db *sql.DB, initDb bool) EmployeeRepository { ... }

// ❌ wrong — exposes concrete type
func NewEmployeeRepository(db *sql.DB) *EmployeeRepositoryImpl { ... }
```
