package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MarkoLuna/AuthService/internal/models"
)

type UserRepositoryImpl struct {
	db          *sql.DB
	useMock     bool
}

func NewUserRepository(db *sql.DB, useMock bool) UserRepository {
	return &UserRepositoryImpl{db: db, useMock: useMock}
}

func (r *UserRepositoryImpl) FindById(id string) (*models.User, error) {
	query := "SELECT id, username, password_hash, email, first_name, last_name, enabled, created_at, updated_at FROM users WHERE id = $1"
	row := r.db.QueryRow(query, id)

	user := &models.User{}
	var email, firstName, lastName sql.NullString
	var createdAt, updatedAt time.Time

	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash,
		&email, &firstName, &lastName,
		&user.Enabled, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	user.Email = email.String
	user.FirstName = firstName.String
	user.LastName = lastName.String
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return user, nil
}

func (r *UserRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	query := "SELECT id, username, password_hash, email, first_name, last_name, enabled, created_at, updated_at FROM users WHERE username = $1"
	row := r.db.QueryRow(query, username)

	user := &models.User{}
	var email, firstName, lastName sql.NullString
	var createdAt, updatedAt time.Time

	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash,
		&email, &firstName, &lastName,
		&user.Enabled, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}

	user.Email = email.String
	user.FirstName = firstName.String
	user.LastName = lastName.String
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return user, nil
}

func (r *UserRepositoryImpl) FindAll() ([]models.User, error) {
	query := "SELECT id, username, password_hash, email, first_name, last_name, enabled, created_at, updated_at FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("find all users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{}
		var email, firstName, lastName sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(&user.Id, &user.Username, &user.PasswordHash,
			&email, &firstName, &lastName,
			&user.Enabled, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		user.Email = email.String
		user.FirstName = firstName.String
		user.LastName = lastName.String
		user.CreatedAt = createdAt
		user.UpdatedAt = updatedAt

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepositoryImpl) Create(user models.User) error {
	query := "INSERT INTO users (id, username, password_hash, email, first_name, last_name, enabled, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err := r.db.Exec(query, user.Id, user.Username, user.PasswordHash,
		nullString(user.Email), nullString(user.FirstName), nullString(user.LastName),
		user.Enabled, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) Update(user models.User) error {
	query := "UPDATE users SET username = $1, password_hash = $2, email = $3, first_name = $4, last_name = $5, enabled = $6, updated_at = $7 WHERE id = $8"
	_, err := r.db.Exec(query, user.Username, user.PasswordHash,
		nullString(user.Email), nullString(user.FirstName), nullString(user.LastName),
		user.Enabled, time.Now(), user.Id)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(id string) error {
	query := "UPDATE users SET enabled = false, updated_at = $1 WHERE id = $2"
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
