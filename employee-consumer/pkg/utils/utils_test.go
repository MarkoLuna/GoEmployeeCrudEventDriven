package utils

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/MarkoLuna/EmployeeConsumer/internal/models"
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

func TestParseIntEnv(t *testing.T) {
	t.Run("DefaultValueWhenUnset", func(t *testing.T) {
		val := ParseIntEnv("UNSET_INT_VAR", 42)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("ValueFromEnv", func(t *testing.T) {
		os.Setenv("SET_INT_VAR", "100")
		defer os.Unsetenv("SET_INT_VAR")
		val := ParseIntEnv("SET_INT_VAR", 42)
		if val != 100 {
			t.Errorf("expected 100, got %d", val)
		}
	})

	t.Run("DefaultValueOnInvalidEnv", func(t *testing.T) {
		os.Setenv("INVALID_INT_VAR", "not-an-int")
		defer os.Unsetenv("INVALID_INT_VAR")
		val := ParseIntEnv("INVALID_INT_VAR", 42)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("DefaultValueOnNonPositiveInt", func(t *testing.T) {
		os.Setenv("ZERO_INT_VAR", "0")
		defer os.Unsetenv("ZERO_INT_VAR")
		val := ParseIntEnv("ZERO_INT_VAR", 42)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})
}

func TestParseBoolEnv(t *testing.T) {
	t.Run("DefaultValueWhenUnset", func(t *testing.T) {
		val := ParseBoolEnv("UNSET_BOOL_VAR", true)
		if val != true {
			t.Errorf("expected true, got %v", val)
		}
	})

	t.Run("ValueFromEnvTrue", func(t *testing.T) {
		os.Setenv("SET_BOOL_VAR", "true")
		defer os.Unsetenv("SET_BOOL_VAR")
		val := ParseBoolEnv("SET_BOOL_VAR", false)
		if val != true {
			t.Errorf("expected true, got %v", val)
		}
	})

	t.Run("ValueFromEnvFalse", func(t *testing.T) {
		os.Setenv("SET_BOOL_VAR", "false")
		defer os.Unsetenv("SET_BOOL_VAR")
		val := ParseBoolEnv("SET_BOOL_VAR", true)
		if val != false {
			t.Errorf("expected false, got %v", val)
		}
	})

	t.Run("DefaultValueOnInvalidEnv", func(t *testing.T) {
		os.Setenv("INVALID_BOOL_VAR", "not-a-bool")
		defer os.Unsetenv("INVALID_BOOL_VAR")
		val := ParseBoolEnv("INVALID_BOOL_VAR", true)
		if val != true {
			t.Errorf("expected true, got %v", val)
		}
	})
}

