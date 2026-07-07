package impl

import (
	"errors"
	"fmt"
	"time"

	"github.com/MarkoLuna/AuthService/internal/models"
	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/MarkoLuna/AuthService/internal/services"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) services.UserService {
	return &UserServiceImpl{userRepo: userRepo}
}

func (s *UserServiceImpl) GetUserId(username string, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", fmt.Errorf("authentication error")
	}
	if user == nil {
		return "", errors.New("invalid username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	return user.Id, nil
}

func (s *UserServiceImpl) CreateUser(req models.UserRequest) (*models.User, error) {
	existing, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	now := time.Now()
	user := models.User{
		Id:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hash),
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Enabled:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserServiceImpl) GetUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}
	return responses, nil
}

func (s *UserServiceImpl) GetUserById(id string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	response := user.ToResponse()
	return &response, nil
}

func (s *UserServiceImpl) UpdateUser(id string, req models.UserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	if user == nil {
		return nil, nil
	}

	user.Username = req.Username
	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.UpdatedAt = time.Now()

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("update user: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	if err := s.userRepo.Update(*user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserServiceImpl) DeleteUser(id string) error {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}
