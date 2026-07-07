package impl

import (
	"errors"

	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/MarkoLuna/AuthService/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type ClientServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewClientService(userRepo repositories.UserRepository) services.ClientService {
	return &ClientServiceImpl{userRepo: userRepo}
}

func (s *ClientServiceImpl) IsValidClientCredentials(clientId string, clientSecret string) (bool, error) {
	client, err := s.userRepo.FindByUsername(clientId)
	if err != nil {
		return false, errors.New("invalid client id")
	}
	if client == nil {
		return false, errors.New("invalid client id")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.PasswordHash), []byte(clientSecret)); err != nil {
		return false, errors.New("invalid client secret")
	}

	return true, nil
}
