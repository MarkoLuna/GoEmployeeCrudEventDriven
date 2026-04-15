package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/MarkoLuna/EmployeeService/internal/models"
	"github.com/MarkoLuna/EmployeeService/internal/services/stubs"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/constants"

	"github.com/stretchr/testify/assert"
)

var (
	basePath = "http://localhost:8080"
)

func InitServer() {
	App.OAuthService = stubs.NewOAuthServiceStub()
	go main()
}

var e = &models.Employee{
	Id:               "1",
	FirstName:        "Marcos",
	LastName:         "Luna",
	SecondLastName:   "Valdez",
	DateOfBirth:      time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC),
	DateOfEmployment: time.Now().UTC(),
	Status:           constants.ACTIVE,
}

func TestMain(m *testing.M) {
	InitServer()

	go func() {
		code := m.Run()
		os.Exit(code)
	}()
}

func makeRequest(httpMethod string, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func TestHealthcheck(t *testing.T) {
	url := fmt.Sprintf("%s/healthcheck/", basePath)
	resp := makeRequest("GET", url, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid http status code")
}
