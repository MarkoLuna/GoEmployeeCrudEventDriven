FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app/auth-service

COPY auth-service/go.mod auth-service/go.sum ./
COPY common/go.mod common/go.sum /app/common/
RUN go mod download

COPY common/ /app/common/
COPY auth-service/ ./

RUN CGO_ENABLED=0 go build -o /dist/auth-service ./cmd/api-server/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates curl

WORKDIR /dist

COPY --from=builder /dist/auth-service .
COPY auth-service/resources ./resources

EXPOSE 8080

CMD ["./auth-service"]
