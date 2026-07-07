package stubs

import (
	"time"

	"github.com/MarkoLuna/AuthService/internal/models"
	"github.com/google/uuid"
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
	now := time.Now()
	user := models.User{
		Id:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.users[user.Id] = user
	return &user, nil
}

func (s UserServiceStub) GetUsers() ([]models.UserResponse, error) {
	responses := make([]models.UserResponse, 0, len(s.users))
	for _, u := range s.users {
		responses = append(responses, u.ToResponse())
	}
	return responses, nil
}

func (s UserServiceStub) GetUserById(id string) (*models.UserResponse, error) {
	u, exists := s.users[id]
	if !exists {
		return nil, nil
	}
	resp := u.ToResponse()
	return &resp, nil
}

func (s UserServiceStub) UpdateUser(id string, req models.UserRequest) (*models.UserResponse, error) {
	u, exists := s.users[id]
	if !exists {
		return nil, nil
	}

	u.Username = req.Username
	u.Email = req.Email
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.UpdatedAt = time.Now()
	s.users[id] = u

	resp := u.ToResponse()
	return &resp, nil
}

func (s UserServiceStub) DeleteUser(id string) error {
	_, exists := s.users[id]
	if !exists {
		return nil
	}
	delete(s.users, id)
	return nil
}
