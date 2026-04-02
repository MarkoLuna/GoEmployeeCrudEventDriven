package utils

import (
	"encoding/json"
	"io"
)

func ParseBody(body io.Reader, x interface{}) {
	if body == nil {
		return
	}

	if data, err := io.ReadAll(body); err == nil {
		if err := json.Unmarshal(data, x); err != nil {
			return
		}
	}
}
