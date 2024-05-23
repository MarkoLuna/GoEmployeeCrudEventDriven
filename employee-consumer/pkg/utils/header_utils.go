package utils

import (
	"net/http"
	"strings"
)

func GetBearerAuth(headers http.Header) (string, bool) {
	return GetAuthHeader(headers, "Bearer ")
}

func GetBasicAuth(headers http.Header) (string, bool) {
	return GetAuthHeader(headers, "Basic ")
}

func GetAuthHeader(headers http.Header, prefix string) (string, bool) {
	auth := headers.Get("Authorization")
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	return token, token != ""
}
