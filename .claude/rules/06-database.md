# Rule: Database & Repository Patterns

> **Applies to**: all `internal/repositories/` code across both services.

## Parameterized Queries — Mandatory

**Never** concatenate user input into SQL strings:

```go
// ✅ safe
err := er.db.QueryRow(
    `SELECT id_employee, first_name, last_name FROM employees WHERE id_employee = $1`, ID,
).Scan(&e.Id, &e.FirstName, &e.LastName)

// ❌ SQL injection risk — never do this
query := "SELECT ... WHERE id_employee = '" + ID + "'"
```

## Always Check `rows.Err()` After Iteration

```go
rows, err := er.db.Query(`SELECT id_employee, first_name FROM employees`)
if err != nil {
    return nil, err
}
defer rows.Close()

var employees []models.Employee
for rows.Next() {
    var e models.Employee
    if err := rows.Scan(&e.Id, &e.FirstName); err != nil {
        log.Printf("scan error: %v", err)
        continue
    }
    employees = append(employees, e)
}
if err := rows.Err(); err != nil {
    return employees, err
}
return employees, nil
```

## Use `RETURNING` for Insert Statements

```go
sqlStatement := `
    INSERT INTO employees (id_employee, first_name, last_name, middle_name, date_of_birth)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id_employee`

err := er.db.QueryRow(sqlStatement, e.Id, e.FirstName, e.LastName, e.MiddleName, e.DateOfBirth).
    Scan(&e.Id)
```

## Interface Contract

Repository interfaces are declared in `internal/repositories/` and returned by factory functions:

```go
// Interface declaration — in repositories package
type EmployeeRepository interface {
    Create(e models.Employee) (*models.Employee, error)
    FindAll() ([]models.Employee, error)
    FindById(ID string) (models.Employee, error)
    DeleteById(ID string) (int64, error)
    Update(e models.Employee) (int64, error)
}

// Factory — returns interface, not concrete type
func NewEmployeeRepository(db *sql.DB, initDb bool) EmployeeRepository {
    repo := &EmployeeRepositoryImpl{db}
    if initDb {
        repo.CreateTable()
    }
    return repo
}
```

## Environment Variables for DB Config

Read via `utils.GetEnv` — never `os.Getenv` directly:

```go
host     := utils.GetEnv("DB_HOST", "localhost")
port     := utils.GetEnv("DB_PORT", "5432")
dbname   := utils.GetEnv("DB_NAME", "employee_db")
user     := utils.GetEnv("DB_USER", "employee_user")
password := utils.GetEnv("DB_PASSWORD", "employeepw")
driver   := utils.GetEnv("DB_DRIVER_NAME", "postgres")
```
