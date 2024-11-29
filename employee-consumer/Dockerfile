FROM golang:1.22-buster

WORKDIR /dist

COPY EmployeeConsumer main
COPY resources/ssl ssl

EXPOSE 8080

CMD ["/dist/main"]

# docker build -t goemployee_crud .