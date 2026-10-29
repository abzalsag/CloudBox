package auth

import (
	"context"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, t *entity.RefreshToken) error
	GetActiveByHash(ctx context.Context, hash string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type TokenManager interface {
	Generate(userID uuid.UUID) (string, error)
}
