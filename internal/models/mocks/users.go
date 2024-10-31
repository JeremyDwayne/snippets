package mocks

import (
	"github.com/jeremydwayne/snippets/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return "e30fd85a-efd2-44d0-86ed-88e71a8dfeda", nil
	}

	return "", models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id string) (bool, error) {
	switch id {
	case "e30fd85a-efd2-44d0-86ed-88e71a8dfeda":
		return true, nil
	default:
		return false, nil
	}
}
