package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

// Insert adds a new user record in the database after hashing the plain text password
func (m *UserModel) Insert(name, email, password string) error {
	stmt := `
	INSERT INTO users(name, email, hashed_password, created)
	VALUES($1, $2, $3, NOW())
	`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(context.Background(), stmt, name, email, string(hashedPassword))
	if err != nil {
		// check if the email address is already in use (by checking for unique constraint violation)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrDuplicateEmail
		}

		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `
	SELECT id, hashed_password
	FROM users
	WHERE email = $1
	`

	id := 0
	hashedPassword := []byte{}
	if err := m.DB.QueryRow(context.Background(), stmt, email).Scan(&id, &hashedPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}

		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}

		return 0, err
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
