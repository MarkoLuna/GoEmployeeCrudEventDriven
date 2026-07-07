package repositories

import "github.com/MarkoLuna/AuthService/internal/models"

type UserRepository interface {
	FindById(id string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindAll() ([]models.User, error)
	Create(user models.User) error
	Update(user models.User) error
	Delete(id string) error
}
