package user

import (
	"context"

	db "CloudBox/db/sqlc"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	email string,
	passwordHash string,
) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
}

func (r *Repository) GetByEmail(
	ctx context.Context,
	email string,
) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}
