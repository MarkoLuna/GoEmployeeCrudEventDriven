package repositories

import (
	"github.com/MarkoLuna/AuthService/internal/models"
)

type UserRepositoryStub struct {
	users map[string]models.User
}

func NewUserRepositoryStub() UserRepository {
	return &UserRepositoryStub{
		users: make(map[string]models.User),
	}
}

func (r *UserRepositoryStub) FindById(id string) (*models.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return &user, nil
}

func (r *UserRepositoryStub) FindByUsername(username string) (*models.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, nil
}

func (r *UserRepositoryStub) FindAll() ([]models.User, error) {
	users := make([]models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryStub) Create(user models.User) error {
	r.users[user.Id] = user
	return nil
}

func (r *UserRepositoryStub) Update(user models.User) error {
	r.users[user.Id] = user
	return nil
}

func (r *UserRepositoryStub) Delete(id string) error {
	if user, exists := r.users[id]; exists {
		user.Enabled = false
		r.users[id] = user
	}
	return nil
}
