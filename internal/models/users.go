package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
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

	query := `INSERT INTO users (id, name, email, hashed_password, created)
            VALUES(?, ?, ?, ?, DATETIME('now', 'utc'))`
	id := uuid.New().String()
	_, err = m.DB.Exec(query, id, name, email, string(hashedPassword))
	if err != nil {
		var sqliteError *sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	var id string
	var hashedPassword []byte

	query := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(query, id).Scan(&exists)
	return exists, err
}
