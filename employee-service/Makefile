# constants
NAME=employee-service
PROJECT?=github.com/MarkoLuna/EmployeeService

.PHONY: help verify build test test-cover vet test-total-cover run clean docker-build docker-run docker-compose-run docker-compose-down k8-apply k8-remove

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## verify: Verify dependencies
verify:
	go mod verify

## build: Build the binary
build: 
	go build -mod readonly -o ${NAME} "${PROJECT}/cmd/api-server"

## test: Run unit tests
test:
	go test -timeout 30s "${PROJECT}/..."

## test-cover: Run unit tests with coverage
test-cover:
	go test -cover "${PROJECT}/..."

## vet: Run go vet
vet:
	go vet "${PROJECT}/..."

## test-total-cover: Run unit tests and show coverage report
test-total-cover:
	go test "${PROJECT}/..." -coverprofile cover.out > /dev/null
	go tool cover -func cover.out
	rm cover.out

## run: Build and run the service
run: build
	./${NAME}

## clean: Clean build artifacts
clean:
	go clean "${PROJECT}/..."
	rm -f ${NAME}

## docker-build: Build docker image
docker-build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${NAME} "${PROJECT}/cmd/api-server"
	docker build -t go-employee-service:latest .
	rm -f ${NAME}

## docker-run: Build and run docker image
docker-run: docker-build
	docker run -it -p 8080:8080 --rm go-employee-service

## docker-compose-run: Build and run with docker-compose
docker-compose-run: docker-build
	docker-compose up

## docker-compose-down: Stop and remove docker-compose containers
docker-compose-down:
	docker-compose down

## k8-apply: Build and apply to Kubernetes
k8-apply: docker-build
	kubectl apply -f ../../k8s/employee-service-pod.yaml
	kubectl apply -f ../../k8s/employee-service-svc.yaml
	kubectl apply -f ../../k8s/ingress.yaml

## k8-remove: Remove from Kubernetes
k8-remove:
	kubectl delete -f ../../k8s/employee-service-pod.yaml
	kubectl delete -f ../../k8s/employee-service-svc.yaml
	kubectl delete -f ../../k8s/ingress.yaml
