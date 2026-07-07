FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y librdkafka-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app/employee-consumer

COPY employee-consumer/go.mod employee-consumer/go.sum ./
COPY common/go.mod common/go.sum /app/common/
RUN go mod download

COPY common/ /app/common/
COPY employee-consumer/ ./

RUN CGO_ENABLED=1 go build -o /dist/employee-consumer ./cmd/api-server/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates curl librdkafka1 && rm -rf /var/lib/apt/lists/*

WORKDIR /dist

COPY --from=builder /dist/employee-consumer .
COPY employee-consumer/resources ./resources

EXPOSE 8080

CMD ["./employee-consumer"]
