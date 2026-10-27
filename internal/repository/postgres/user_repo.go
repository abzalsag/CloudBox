package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"CloudBox/internal/entity"
	"CloudBox/internal/repository/postgres/sqlcgen"
)

type UserRepo struct {
	*Store
}

func NewUserRepo(s *Store) *UserRepo { return &UserRepo{Store: s} }

func (r *UserRepo) Create(ctx context.Context, u *entity.User) error {
	row, err := r.Q(ctx).CreateUser(ctx, sqlcgen.CreateUserParams{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return entity.ErrAlreadyExists
		}
		return err
	}
	*u = *toUser(row)
	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row, err := r.Q(ctx).GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toUser(row), nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row, err := r.Q(ctx).GetUserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toUser(row), nil
}

// IncreaseStorageUsed: квота асса entity.ErrQuotaExceeded қайтарады.
func (r *UserRepo) IncreaseStorageUsed(ctx context.Context, id uuid.UUID, delta int64) error {
	_, err := r.Q(ctx).IncreaseStorageUsed(ctx, sqlcgen.IncreaseStorageUsedParams{
		ID:    id,
		Delta: delta,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.ErrQuotaExceeded
	}
	return err
}

func (r *UserRepo) DecreaseStorageUsed(ctx context.Context, id uuid.UUID, delta int64) error {
	return r.Q(ctx).DecreaseStorageUsed(ctx, sqlcgen.DecreaseStorageUsedParams{
		ID:    id,
		Delta: delta,
	})
}

func toUser(row sqlcgen.User) *entity.User {
	return &entity.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName,
		StorageQuota: row.StorageQuota,
		StorageUsed:  row.StorageUsed,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
