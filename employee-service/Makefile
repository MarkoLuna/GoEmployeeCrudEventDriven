# constants
NAME=employee-service
PROJECT?=github.com/MarkoLuna/EmployeeService

verify:
	go mod verify

build: 
	go build -mod readonly -o ${NAME} "${PROJECT}/cmd/api-server"

test:
	go test -timeout 30s "${PROJECT}/..."

test-cover:
	go test -cover "${PROJECT}/..."

vet:
	go vet "${PROJECT}/..."

test-total-cover:
	go test "${PROJECT}/..." -coverprofile cover.out > /dev/null
	go tool cover -func cover.out
	rm cover.out

run: build
	./${NAME}

clean:
	go clean "${PROJECT}/..."
	rm -f ${NAME}

docker-build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${NAME} "${PROJECT}/cmd/api-server"
	docker build -t go-employee-service:latest .
	rm -f ${NAME}

docker-run: docker-build
	docker run -it -p 8080:8080 --rm go-employee-service

docker-compose-run: docker-build
	docker-compose up

docker-compose-down:
	docker-compose down

k8-apply: docker-build
	kubectl apply -f ../../k8s/employee-service-pod.yaml
	kubectl apply -f ../../k8s/employee-service-svc.yaml
	kubectl apply -f ../../k8s/ingress.yaml

k8-remove:
	kubectl delete -f ../../k8s/employee-service-pod.yaml
	kubectl delete -f ../../k8s/employee-service-svc.yaml
	kubectl delete -f ../../k8s/ingress.yaml
