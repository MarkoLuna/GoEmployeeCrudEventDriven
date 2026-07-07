package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_ToResponse(t *testing.T) {
	now := time.Now()
	client := Client{
		Id:           "client-1",
		ClientId:     "my-client",
		ClientSecret: "should-not-appear",
		Enabled:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	resp := client.ToResponse()

	assert.Equal(t, "client-1", resp.Id)
	assert.Equal(t, "my-client", resp.ClientId)
	assert.True(t, resp.Enabled)
	assert.Equal(t, now, resp.CreatedAt)
	assert.Equal(t, now, resp.UpdatedAt)
}

func TestClient_ToResponse_Empty(t *testing.T) {
	client := Client{}

	resp := client.ToResponse()

	assert.Empty(t, resp.Id)
	assert.Empty(t, resp.ClientId)
	assert.False(t, resp.Enabled)
}
