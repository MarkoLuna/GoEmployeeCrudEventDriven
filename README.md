# Employee Crud

Employee Crud Rest API using Golang and Event Driven Architecture

## Install deps

```bash
$ go mod tidy
```

## Run on local

```
$ go run pkg/main.go
```

Or with make

```
$ make run
```

## Run with docker

```
$ make docker-run
```

## Run with docker compose

```
$ make docker-compose-run
```

## healthcheck

```
$ curl -X GET http://localhost:8080/healthcheck/
OK
```

## Run in k8 locally

```
$ make k8-apply
$ kubectl exec -it employeecrud-pod -- apk add curl
$ kubectl exec -it employeecrud-pod -- curl -X GET http://employeecrud-service:8080/healthcheck/
$ make k8-remove 
```

## Example Curls

```
# get all employees
$ curl http://localhost:8080/api/employee/
[]

# get employee by id
$ curl http://localhost:8080/api/employee/1 

# create employee
$ curl --location --request POST 'http://localhost:8080/api/employee/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "firstName": "Marcos",
    "lastName": "Luna",
    "secondLastName": "Valdez",
    "dateOfBirth": "1994-04-25T12:00:00Z",
    "dateOfEmployment": "1994-04-25T12:00:00Z",
    "status": "ACTIVE"
}'

# delete employee
$ curl -X DELETE 'http://localhost:8080/api/employee/2'

# update employee
$ curl --location --request PUT 'http://localhost:8080/api/employee/3' \
--header 'Content-Type: application/json' \
--data-raw '{
    "firstName": "Gerardo",
    "lastName": "Luna",
    "secondLastName": "Valdezz",
    "dateOfBirth": "1994-04-25T12:00:00Z",
    "dateOfEmployment": "0001-01-01T00:00:00Z",
    "status": "INACTIVE"
}'
```
