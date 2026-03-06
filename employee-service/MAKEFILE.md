# Makefile Documentation — `employee-service`

## Variables

| Variable  | Default Value                              | Description                                                |
|-----------|--------------------------------------------|------------------------------------------------------------|
| `NAME`    | `employee-service`                         | Name of the compiled binary output.                        |
| `PROJECT` | `github.com/MarkoLuna/EmployeeService`     | Go module path. Can be overridden via environment variable.|

---

## Targets

### Development

| Target         | Command(s)                                         | Description                                                                 |
|----------------|----------------------------------------------------|-----------------------------------------------------------------------------|
| `verify`       | `go mod verify`                                    | Verifies that module dependencies match their checksums.                    |
| `build`        | `go build -mod readonly -o ${NAME} ...`            | Compiles the `cmd/api-server` entrypoint into the `employeeCrudApp` binary. |
| `run`          | Depends on `build`, then `./${NAME}`               | Builds and immediately runs the binary.                                     |
| `clean`        | `go clean ...` + `rm -f ${NAME}`                   | Removes Go build cache and the compiled binary.                             |
| `vet`          | `go vet ...`                                       | Runs static analysis to catch common Go mistakes.                           |

### Testing

| Target            | Command(s)                                          | Description                                                              |
|-------------------|-----------------------------------------------------|--------------------------------------------------------------------------|
| `test`            | `go test -timeout 30s ...`                          | Runs all tests with a 30-second timeout.                                 |
| `test-cover`      | `go test -cover ...`                                | Runs all tests and prints a basic coverage summary.                      |
| `test-total-cover`| `go test ... -coverprofile` + `go tool cover -func` | Generates a full per-function coverage report, then removes the profile. |

### Docker

| Target                | Command(s)                                                       | Description                                                                         |
|-----------------------|------------------------------------------------------------------|-------------------------------------------------------------------------------------|
| `docker-build`        | Cross-compile for Linux/amd64 + `docker build`                  | Builds a Linux binary and packages it into the `goemployee_service:latest` image.   |
| `docker-run`          | Depends on `docker-build`, then `docker run`                     | Builds the image and runs the container, exposing port `8080`.                      |
| `docker-compose-run`  | Depends on `docker-build`, then `docker-compose up`              | Builds the image and starts the full stack via Docker Compose.                      |
| `docker-compose-down` | `docker-compose down`                                            | Stops and removes containers started by Docker Compose.                             |

### Kubernetes

| Target      | Command(s)                                    | Description                                                          |
|-------------|-----------------------------------------------|----------------------------------------------------------------------|
| `k8-apply`  | Depends on `docker-build`, then `kubectl apply`| Builds the image and deploys the pod and service from `k8s/`.        |
| `k8-remove` | `kubectl delete pod/service ...`              | Removes the `employeeservice-pod` pod and `employeeservice-service`. |

---

## Common Usage

```bash
# Build and run locally
make run

# Run all tests with coverage report
make test-total-cover

# Build Docker image and start with Compose
make docker-compose-run

# Stop Docker Compose stack
make docker-compose-down

# Deploy to Kubernetes
make k8-apply

# Remove Kubernetes resources
make k8-remove

# Clean build artifacts
make clean
```
