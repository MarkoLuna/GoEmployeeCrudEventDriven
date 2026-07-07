package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_ToResponse(t *testing.T) {
	now := time.Now()
	user := User{
		Id:           "user-1",
		Username:     "testuser",
		PasswordHash: "should-not-appear",
		Email:        "test@test.com",
		FirstName:    "Test",
		LastName:     "User",
		Enabled:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	resp := user.ToResponse()

	assert.Equal(t, "user-1", resp.Id)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@test.com", resp.Email)
	assert.Equal(t, "Test", resp.FirstName)
	assert.Equal(t, "User", resp.LastName)
	assert.True(t, resp.Enabled)
	assert.Equal(t, now, resp.CreatedAt)
	assert.Equal(t, now, resp.UpdatedAt)
}

func TestUser_ToResponse_Empty(t *testing.T) {
	user := User{}

	resp := user.ToResponse()

	assert.Empty(t, resp.Id)
	assert.Empty(t, resp.Username)
	assert.Empty(t, resp.Email)
	assert.False(t, resp.Enabled)
}
