package utils

import (
	"net/http"
	"testing"
)

func TestCannotGetBearerHeader(t *testing.T) {

	headers := http.Header{}
	_, ok := GetBearerAuth(headers)

	if ok {
		t.Error("should not contain the header authorization")
	}
}

func TestGetBearerHeader(t *testing.T) {

	headers := http.Header{}
	headers.Set("Authorization", "Bearer abc")

	token, ok := GetBearerAuth(headers)

	if !ok {
		t.Error("should contain the authorization header")
	}

	if token != "abc" {
		t.Error("incorrect bearer authorization value")
	}
}

func TestCannotGetBasicHeader(t *testing.T) {

	headers := http.Header{}
	_, ok := GetBasicAuth(headers)

	if ok {
		t.Error("should not contain the authorization header")
	}
}

func TestGetBasicHeader(t *testing.T) {

	headers := http.Header{}
	headers.Set("Authorization", "Basic abc")

	token, ok := GetBasicAuth(headers)

	if !ok {
		t.Error("should contain the authorization header")
	}

	if token != "abc" {
		t.Error("incorrect basic authorization value")
	}
}
