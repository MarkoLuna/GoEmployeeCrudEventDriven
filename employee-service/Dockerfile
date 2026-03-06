FROM golang:1.25-alpine

WORKDIR /dist

COPY employee-service main
COPY resources/ssl ssl

EXPOSE 8080

CMD ["/dist/main"]

# docker build -t go-employee-service .