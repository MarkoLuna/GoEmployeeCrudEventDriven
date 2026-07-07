.PHONY: run build test test-cover clean

APP_NAME = auth-service
CMD_DIR = ./cmd/api-server

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

run: build
	./bin/$(APP_NAME)

test:
	go test ./internal/...

test-cover:
	go test ./internal/... -coverprofile=coverage.out -covermode=atomic

vet:
	go vet ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t $(APP_NAME):latest .

tidy:
	go mod tidy
