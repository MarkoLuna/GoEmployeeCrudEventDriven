package impl

import (
	"testing"
	"time"

	"github.com/MarkoLuna/AuthService/internal/models"
	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	return string(hash)
}

func preloadUser(repo repositories.UserRepository, user models.User) {
	repo.Create(user)
}

func TestUserServiceImpl_GetUserId_Success(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	password := "secret123"
	preloadUser(repo, models.User{
		Id:           "user-1",
		Username:     "testuser",
		PasswordHash: hashPassword(t, password),
		Enabled:      true,
	})

	userId, err := svc.GetUserId("testuser", password)
	assert.NoError(t, err)
	assert.Equal(t, "user-1", userId)
}

func TestUserServiceImpl_GetUserId_InvalidUsername(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	_, err := svc.GetUserId("nonexistent", "whatever")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid username")
}

func TestUserServiceImpl_GetUserId_InvalidPassword(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	preloadUser(repo, models.User{
		Id:           "user-1",
		Username:     "testuser",
		PasswordHash: hashPassword(t, "correct-password"),
		Enabled:      true,
	})

	_, err := svc.GetUserId("testuser", "wrong-password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
}

func TestUserServiceImpl_CreateUser_Success(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	req := models.UserRequest{
		Username:  "newuser",
		Password:  "pass123",
		Email:     "new@test.com",
		FirstName: "New",
		LastName:  "User",
	}

	user, err := svc.CreateUser(req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	assert.Equal(t, "new@test.com", user.Email)
	assert.Equal(t, "New", user.FirstName)
	assert.Equal(t, "User", user.LastName)
	assert.True(t, user.Enabled)
	assert.NotEmpty(t, user.Id)
	assert.NotEmpty(t, user.PasswordHash)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("pass123"))
	assert.NoError(t, err, "password should be bcrypt-hashed")
}

func TestUserServiceImpl_CreateUser_DuplicateUsername(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	preloadUser(repo, models.User{
		Id:       "existing-id",
		Username: "duplicate",
	})

	req := models.UserRequest{
		Username: "duplicate",
		Password: "pass123",
	}

	user, err := svc.CreateUser(req)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "username already exists")
}

func TestUserServiceImpl_GetUsers_Success(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	preloadUser(repo, models.User{
		Id:       "1",
		Username: "alice",
		Email:    "alice@test.com",
		Enabled:  true,
	})
	preloadUser(repo, models.User{
		Id:       "2",
		Username: "bob",
		Email:    "bob@test.com",
		Enabled:  true,
	})

	users, err := svc.GetUsers()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
}

func TestUserServiceImpl_GetUsers_Empty(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	users, err := svc.GetUsers()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(users))
}

func TestUserServiceImpl_GetUserById_Found(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	preloadUser(repo, models.User{
		Id:       "user-1",
		Username: "testuser",
		Email:    "test@test.com",
		Enabled:  true,
	})

	user, err := svc.GetUserById("user-1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
}

func TestUserServiceImpl_GetUserById_NotFound(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	user, err := svc.GetUserById("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserServiceImpl_UpdateUser_Success(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	preloadUser(repo, models.User{
		Id:       "user-1",
		Username: "oldname",
		Email:    "old@test.com",
		Enabled:  true,
		CreatedAt: time.Now(),
	})

	req := models.UserRequest{
		Username:  "newname",
		Email:     "new@test.com",
		FirstName: "Updated",
		LastName:  "User",
		Password:  "newpass",
	}

	user, err := svc.UpdateUser("user-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newname", user.Username)
	assert.Equal(t, "new@test.com", user.Email)
	assert.Equal(t, "Updated", user.FirstName)
}

func TestUserServiceImpl_UpdateUser_NotFound(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	req := models.UserRequest{
		Username: "any",
		Password: "pass",
	}

	user, err := svc.UpdateUser("nonexistent", req)
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserServiceImpl_DeleteUser_Success(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	preloadUser(repo, models.User{
		Id:       "user-1",
		Username: "todelete",
		Enabled:  true,
	})

	err := svc.DeleteUser("user-1")
	assert.NoError(t, err)
}

func TestUserServiceImpl_DeleteUser_NotFound(t *testing.T) {
	repo := repositories.NewUserRepositoryStub()
	svc := NewUserService(repo)

	err := svc.DeleteUser("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}
