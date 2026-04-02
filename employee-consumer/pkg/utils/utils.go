package utils

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
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

func ParseIntEnv(key string, defaultVal int) int {
	raw := GetEnv(key, strconv.Itoa(defaultVal))
	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return defaultVal
	}
	return val
}

func ParseBoolEnv(key string, defaultVal bool) bool {
	raw := GetEnv(key, strconv.FormatBool(defaultVal))
	val, err := strconv.ParseBool(raw)
	if err != nil {
		return defaultVal
	}
	return val
}
