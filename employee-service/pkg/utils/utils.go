package utils

import (
	"encoding/json"
	"io"
	"os"
)

func ParseBody(body io.Reader, x interface{}) {

	if body == nil {
		return
	}

	if body, err := io.ReadAll(body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
