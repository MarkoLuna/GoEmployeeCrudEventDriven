package stubs

import (
	"github.com/MarkoLuna/AuthService/internal/models"
)

type UserServiceStub struct {
	users map[string]models.User
}

func NewUserServiceStub() UserServiceStub {
	return UserServiceStub{
		users: make(map[string]models.User),
	}
}

func (s UserServiceStub) GetUserId(username string, password string) (string, error) {
	for _, u := range s.users {
		if u.Username == username {
			return u.Id, nil
		}
	}
	return "", nil
}

func (s UserServiceStub) CreateUser(req models.UserRequest) (*models.User, error) {
	return nil, nil
}

func (s UserServiceStub) GetUsers() ([]models.UserResponse, error) {
	return nil, nil
}

func (s UserServiceStub) GetUserById(id string) (*models.UserResponse, error) {
	return nil, nil
}

func (s UserServiceStub) UpdateUser(id string, req models.UserRequest) (*models.UserResponse, error) {
	return nil, nil
}

func (s UserServiceStub) DeleteUser(id string) error {
	return nil
}
