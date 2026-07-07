package services

import (
	"github.com/MarkoLuna/AuthService/internal/models"
)

type UserService interface {
	GetUserId(username string, password string) (string, error)
	CreateUser(req models.UserRequest) (*models.User, error)
	GetUsers() ([]models.UserResponse, error)
	GetUserById(id string) (*models.UserResponse, error)
	UpdateUser(id string, req models.UserRequest) (*models.UserResponse, error)
	DeleteUser(id string) error
}
