FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y librdkafka-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app/employee-service

COPY employee-service/go.mod employee-service/go.sum ./
COPY common/go.mod common/go.sum /app/common/
RUN go mod download

COPY common/ /app/common/
COPY employee-service/ ./

RUN CGO_ENABLED=1 go build -o /dist/employee-service ./cmd/api-server/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates curl librdkafka1 && rm -rf /var/lib/apt/lists/*

WORKDIR /dist

COPY --from=builder /dist/employee-service .
COPY employee-service/resources ./resources

EXPOSE 8080

CMD ["./employee-service"]
