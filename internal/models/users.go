package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jeremydwayne/snippets/internal/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"hashed_password"`
	Created        time.Time `json:"created"`
}

type UserModel struct {
	DB *sqlc.Queries
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (string, error)
	Exists(id string) (bool, error)
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	id := uuid.New().String()
	params := sqlc.CreateUserParams{
		ID:             id,
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	err = m.DB.CreateUser(context.Background(), params)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	var user sqlc.AuthenticateUserRow

	user, err := m.DB.AuthenticateUser(context.Background(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	return user.ID, nil
}

func (m *UserModel) Exists(id string) (bool, error) {
	exists, err := m.DB.UserExists(context.Background(), id)
	return exists != 0, err
}
