package utils

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/MarkoLuna/EmployeeService/pkg/models"
)

func TestGetEnvDefaultValue(t *testing.T) {

	DbHost := GetEnv("DB_HOST", "localhost")
	if DbHost != "localhost" {
		t.Error("should be the default value")
	}
}

func TestGetEnvFromSO(t *testing.T) {

	os.Setenv("DB_HOST", "localhost1")

	DbHost := GetEnv("DB_HOST", "localhost")
	if DbHost != "localhost1" {
		t.Error("should be the env variable value")
	}

	os.Unsetenv("DB_HOST")
}

func TestParseBodyWithNil(t *testing.T) {

	req, err := http.NewRequest("GET", "/healthcheck/", nil)
	if err != nil {
		t.Fatal(err)
	}

	employee := &models.Employee{}
	ParseBody(req.Body, employee)
	// fmt.Println("employee: " + employee.ToString())

	v := CreateValidator()
	err1 := v.Struct(employee)
	if err1 == nil {
		t.Error("the body should not be valid")
	}
}

func TestParseBodyWithInvalidData(t *testing.T) {

	var jsonStr = []byte(`
			{
				"client_id":"Marcos", 
				"grant_type": "client_credentials"
			}
			`)

	body := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", "/healthcheck/", body)

	if err != nil {
		t.Fatal(err)
	}

	employee := &models.Employee{}
	ParseBody(req.Body, employee)
	// fmt.Println("employee: " + employee.ToString())

	v := CreateValidator()
	err1 := v.Struct(employee)
	if err1 == nil {
		t.Error("the body should not be valid")
	}
}

func TestParseBodyWithValidData(t *testing.T) {

	var jsonStr = []byte(`
			{
				"firstName":"Marcos", 
				"lastName": "Luna", 
				"secondLastName": "Valdez",
				"dateOfBirth": "1994-04-25T12:00:00Z", 
				"dateOfEmployment": "1994-04-25T12:00:00Z",
				"status": "ACTIVE" 
			}
			`)

	body := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", "/healthcheck/", body)

	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	employee := &models.Employee{}
	ParseBody(req.Body, employee)
	// fmt.Println("employee: " + employee.ToString())

	v := CreateValidator()
	err2 := v.Struct(employee)
	if err2 != nil {
		t.Error("the body should be valid")
	}

}
