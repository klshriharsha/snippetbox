package mocks

import (
	"time"

	"github.com/klshriharsha/snippetbox/internal/models"
)

var mockUser = &models.User{
	ID:             1,
	Name:           "Snippetbox",
	Email:          "snippetbox@example.com",
	HashedPassword: []byte("password"),
	Created:        time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "duplicate@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "snippetbox@example.com" && password == "password" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) GetByID(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
