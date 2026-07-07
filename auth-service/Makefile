# constants
NAME=auth-service
PROJECT?=github.com/MarkoLuna/AuthService

.PHONY: help verify build test test-cover vet test-total-cover run clean docker-build docker-run k8-apply k8-remove

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
	docker build -t go-auth-service:latest .
	rm -f ${NAME}

## docker-run: Build and run docker image
docker-run: docker-build
	docker run -it -p 8082:8082 --rm go-auth-service

## k8-apply: Build and apply to Kubernetes
k8-apply: docker-build
	kubectl apply -f ../k8s/configmap.yaml
	kubectl apply -f ../k8s/secrets.yaml
	kubectl apply -f ../k8s/auth-service-deployment.yaml
	kubectl apply -f ../k8s/auth-service-svc.yaml
	kubectl apply -f ../k8s/ingress.yaml

## k8-remove: Remove from Kubernetes
k8-remove:
	kubectl delete -f ../k8s/auth-service-deployment.yaml
	kubectl delete -f ../k8s/auth-service-svc.yaml
	kubectl delete -f ../k8s/ingress.yaml
	kubectl delete -f ../k8s/configmap.yaml
	kubectl delete -f ../k8s/secrets.yaml
