package api

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mattrax/Mattrax/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// Service abstracts resources to prevent mishandling.
type Service struct {
	DB *db.Queries
}

var ErrIncorrectCredentials = fmt.Errorf("error authentication: invalid credentials")
var ErrUserIsDisabled = fmt.Errorf("error authentication: user login disabled")

// Login authenticates a user to the API
func (s *Service) Login(ctx context.Context, username, password string) (db.GetUserSecureRow, error) {
	user, err := s.DB.GetUserSecure(ctx, username)
	if err == sql.ErrNoRows {
		return db.GetUserSecureRow{}, ErrIncorrectCredentials

	} else if err != nil {
		return db.GetUserSecureRow{}, fmt.Errorf("error retrieving user from DB: %w", err)
	}

	if !user.Password.Valid {
		return db.GetUserSecureRow{}, ErrIncorrectCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password)); err == bcrypt.ErrMismatchedHashAndPassword {
		return db.GetUserSecureRow{}, ErrIncorrectCredentials
	} else if err != nil {
		return db.GetUserSecureRow{}, fmt.Errorf("error comparing password to hash: %w", err)
	}

	if user.Disabled {
		return db.GetUserSecureRow{}, ErrUserIsDisabled
	}

	return user, nil
}

// New initialises a new API service
func New(db *db.Queries) (s *Service, err error) {
	return &Service{
		DB: db,
	}, nil
}
