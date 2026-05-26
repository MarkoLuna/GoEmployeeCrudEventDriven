# Rule: Testing Standards

> **Applies to**: all `*_test.go` files across all modules.

## File Placement

Co-locate the test file with the file under test:

```
internal/services/impl/employee_service_impl.go
internal/services/impl/employee_service_impl_test.go  ← same directory
```

## Stub-Based Unit Testing (Preferred)

Implement the service/repository interface with a hardcoded stub — **no mocking frameworks**:

```go
// internal/services/stubs/kafka_producer_service_stub.go
type KafkaProducerServiceStub struct {
    SendInsertCalled bool
    SendInsertErr    error
}

func (s *KafkaProducerServiceStub) SendInsert(e dto.EmployeeMessage) error {
    s.SendInsertCalled = true
    return s.SendInsertErr
}
```

Stubs live in:
- `internal/services/stubs/` — service stubs
- `internal/repositories/` — repository stubs (alongside interface file)

## Repository Tests with `go-sqlmock`

```go
db, mock, _ := sqlmock.New()
defer db.Close()

mock.ExpectQuery(`SELECT .+ FROM employees`).
    WillReturnRows(sqlmock.NewRows([]string{"id_employee", "first_name", ...}).
        AddRow("1", "John", ...))

repo := NewEmployeeRepository(db, false)
employees, err := repo.FindAll()
assert.NoError(t, err)
assert.Len(t, employees, 1)
assert.NoError(t, mock.ExpectationsWereMet())
```

## HTTP Handler Tests with Echo

```go
e := echo.New()
req := httptest.NewRequest(http.MethodPost, "/api/employee/", strings.NewReader(body))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)

ctrl := controllers.NewEmployeeController(&stubs.EmployeeServiceStub{})
err := ctrl.CreateEmployee(c)
assert.NoError(t, err)
assert.Equal(t, http.StatusCreated, rec.Code)
```

## Running Tests

```bash
cd employee-service && make test
cd employee-consumer && make test
```

## Rules

- Unit tests must **never** connect to a real DB or real Kafka broker.
- Use `testify/assert` for assertions — no raw `t.Errorf` unless necessary.
- Always call `mock.ExpectationsWereMet()` at the end of sqlmock tests.
- Use `assert.NoError(t, err)` before asserting on the result value.
