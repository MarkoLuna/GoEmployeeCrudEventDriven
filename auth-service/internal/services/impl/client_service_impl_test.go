package impl

import (
	"testing"

	"github.com/MarkoLuna/AuthService/internal/models"
	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestClientServiceImpl_IsValidClientCredentials_Valid(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewClientService(repo)

	secret := "client-secret"
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	assert.NoError(t, err)

	preloadUser(repo, models.User{
		Id:           "client-1",
		Username:     "myclient",
		PasswordHash: string(hash),
		Enabled:      true,
	})

	valid, err := svc.IsValidClientCredentials("myclient", secret)
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestClientServiceImpl_IsValidClientCredentials_InvalidClientId(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewClientService(repo)

	valid, err := svc.IsValidClientCredentials("nonexistent", "whatever")
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "invalid client id")
}

func TestClientServiceImpl_IsValidClientCredentials_InvalidSecret(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewClientService(repo)

	secret := "correct-secret"
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	assert.NoError(t, err)

	preloadUser(repo, models.User{
		Id:           "client-1",
		Username:     "myclient",
		PasswordHash: string(hash),
		Enabled:      true,
	})

	valid, err := svc.IsValidClientCredentials("myclient", "wrong-secret")
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "invalid client secret")
}
