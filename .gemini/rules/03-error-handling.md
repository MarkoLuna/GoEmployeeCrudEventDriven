# Rule: Error Handling

> **Applies to**: all Go service, repository, and controller code.

## Principles

1. **Propagate, don't swallow** — return errors to the caller; handle at the boundary (controller).
2. **Wrap with context** — use `fmt.Errorf("operation description: %w", err)` so the full call chain is visible.
3. **Never `panic`** in service or repository code — always return an error.
4. **`log.Fatal` only in `main` or `app.Run`** — never inside library-level functions.

## Patterns

### Repository: return and propagate

```go
func (er EmployeeRepositoryImpl) FindById(ID string) (models.Employee, error) {
    var e models.Employee
    err := er.db.QueryRow(sql, ID).Scan(&e.Id, &e.FirstName ...)
    if err != nil {
        return e, err  // caller decides
    }
    return e, nil
}
```

### Service: wrap with context

```go
if err := kSrv.consumer.SubscribeTopics(topics, nil); err != nil {
    return fmt.Errorf("failed to subscribe to topics: %w", err)
}
```

```go
value, err := json.Marshal(employee)
if err != nil {
    return fmt.Errorf("marshal employee message: %w", err)  // ← not panic
}
```

### Controller: HTTP status mapping

| Condition | Status |
|---|---|
| Bad request body / validation failure | `400 Bad Request` |
| Resource not found | `404 Not Found` |
| Service / DB error | `500 Internal Server Error` |
| Successful creation | `201 Created` |
| Successful read/update/delete | `200 OK` |

```go
func (eCtrl EmployeeController) CreateEmployee(c echo.Context) error {
    req := dto.EmployeeRequest{}
    if err := c.Bind(&req); err != nil {
        return c.String(http.StatusBadRequest, err.Error())
    }
    if err := utils.CreateValidator().Struct(req); err != nil {
        return c.String(http.StatusBadRequest, err.Error())
    }
    result, err := eCtrl.employeeService.CreateEmployee(req)
    if err != nil {
        return c.String(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusCreated, result)
}
```

## Anti-patterns to Avoid

| ❌ Anti-pattern | ✅ Correct approach |
|---|---|
| `panic(err)` on serialisation failure | `return fmt.Errorf("marshal: %w", err)` |
| Returning `nil` error when one occurred | Always propagate the error |
| Logging the error AND returning it | Choose one; prefer returning |
| `log.Fatal` inside `internal/` packages | Only in `main()` or `app.Run()` |
