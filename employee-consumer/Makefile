# constants
NAME=employeeCrudApp
PROJECT?=github.com/MarkoLuna/EmployeeConsumer

# https://github.com/swaggo/echo-swagger
# https://github.com/swaggo/swag#declarative-comments-format
swagger-build:
	swag init --dir pkg

verify:
	go mod verify

build: swagger-build
	go build -mod readonly -o ${NAME} "${PROJECT}/pkg"

test:
	go test -race "${PROJECT}/..."

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
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${NAME} "${PROJECT}/pkg"
	docker build -t goemployee_consumer:latest .
	rm -f ${NAME}

docker-run: docker-build
	docker run -it -p 8080:8080 --rm goemployee_consumer

docker-compose-run: docker-build
	docker-compose up

docker-compose-down:
	docker-compose down

k8-apply: docker-build
	kubectl apply -f k8s/pod.yaml
	kubectl apply -f k8s/service.yaml

k8-remove:
	kubectl delete pod employeeconsumer-pod
	kubectl delete service employeeconsumer-service
