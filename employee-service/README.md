# Employee Crud

Employee Crud Rest and Kafka Producer API using Golang

## Run with docker

```bash
make docker-run
```

## Run with docker compose

```bash
make docker-compose-run
```

## Run in k8 locally

```bash
$ make k8-apply
$ kubectl exec -it employeecrud-service-pod -- apk add curl
$ kubectl exec -it employeecrud-service-pod -- curl -X GET http://employeecrud-service-pod:8080/healthcheck/
$ make k8-remove 
```

## Example Curls

```bash
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
