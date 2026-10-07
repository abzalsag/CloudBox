package user

import (
	"context"
	"errors"

	db "CloudBox/db/sqlc"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) (db.User, error) {
	if email == "" || password == "" {
		return db.User{}, errors.New("email and password are required")
	}

	_, err := s.repository.GetByEmail(ctx, email)
	if err == nil {
		return db.User{}, errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return db.User{}, err
	}

	return s.repository.Create(
		ctx,
		email,
		string(passwordHash),
	)
}
