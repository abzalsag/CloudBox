package postgres

import (
	"context"
	"errors"
	"net/netip"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"CloudBox/internal/entity"
	"CloudBox/internal/repository/postgres/sqlcgen"
)

type RefreshTokenRepo struct {
	*Store
}

func NewRefreshTokenRepo(s *Store) *RefreshTokenRepo { return &RefreshTokenRepo{Store: s} }

func (r *RefreshTokenRepo) Create(ctx context.Context, t *entity.RefreshToken) error {
	var ip *netip.Addr
	if addr, err := netip.ParseAddr(t.IP); err == nil {
		ip = &addr
	}
	_, err := r.Q(ctx).CreateRefreshToken(ctx, sqlcgen.CreateRefreshTokenParams{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenHash: t.TokenHash,
		UserAgent: t.UserAgent,
		Ip:        ip,
		ExpiresAt: t.ExpiresAt,
	})
	return err
}

func (r *RefreshTokenRepo) GetActiveByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	row, err := r.Q(ctx).GetActiveRefreshTokenByHash(ctx, hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	ip := ""
	if row.Ip != nil {
		ip = row.Ip.String()
	}
	return &entity.RefreshToken{
		ID:        row.ID,
		UserID:    row.UserID,
		TokenHash: row.TokenHash,
		UserAgent: row.UserAgent,
		IP:        ip,
		ExpiresAt: row.ExpiresAt,
		RevokedAt: row.RevokedAt,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *RefreshTokenRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.Q(ctx).RevokeRefreshToken(ctx, id)
}

func (r *RefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	return r.Q(ctx).RevokeAllUserRefreshTokens(ctx, userID)
}
