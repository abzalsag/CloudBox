package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type Meta struct {
	UserAgent string
	IP        string
}

type RegisterInput struct {
	Email    string
	Password string
	FullName string
}

type LoginInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // секунд
}

type Usecase struct {
	users      UserRepository
	tokens     RefreshTokenRepository
	tx         TxManager
	hasher     PasswordHasher
	jwt        TokenManager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewUsecase(
	users UserRepository,
	tokens RefreshTokenRepository,
	tx TxManager,
	hasher PasswordHasher,
	jwt TokenManager,
	accessTTL, refreshTTL time.Duration,
) *Usecase {
	return &Usecase{
		users: users, tokens: tokens, tx: tx, hasher: hasher, jwt: jwt,
		accessTTL: accessTTL, refreshTTL: refreshTTL,
	}
}

func (u *Usecase) Register(ctx context.Context, in RegisterInput) (*entity.User, error) {
	email := normalizeEmail(in.Email)
	if addr, err := mail.ParseAddress(email); err != nil || addr.Address != email {
		return nil, fmt.Errorf("%w: invalid email", entity.ErrInvalidInput)
	}
	name := strings.TrimSpace(in.FullName)
	if name == "" {
		return nil, fmt.Errorf("%w: full name is required", entity.ErrInvalidInput)
	}
	// bcrypt 72 байттан артығын ескермейді
	if len(in.Password) < 8 || len(in.Password) > 72 {
		return nil, fmt.Errorf("%w: password must be 8-72 characters", entity.ErrInvalidInput)
	}

	hash, err := u.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	user := &entity.User{ID: id, Email: email, PasswordHash: hash, FullName: name}
	if err := u.users.Create(ctx, user); err != nil {
		return nil, err // ErrAlreadyExists-ті репозиторий қайтарады
	}
	return user, nil
}

func (u *Usecase) Login(ctx context.Context, in LoginInput, meta Meta) (*Tokens, error) {
	user, err := u.users.GetByEmail(ctx, normalizeEmail(in.Email))
	if errors.Is(err, entity.ErrNotFound) {
		return nil, entity.ErrUnauthorized // email бар-жоғын жасырамыз
	}
	if err != nil {
		return nil, err
	}
	if !u.hasher.Compare(user.PasswordHash, in.Password) {
		return nil, entity.ErrUnauthorized
	}
	return u.issueTokens(ctx, user.ID, meta)
}

func (u *Usecase) Refresh(ctx context.Context, rawRefresh string, meta Meta) (*Tokens, error) {
	rt, err := u.tokens.GetActiveByHash(ctx, hashToken(rawRefresh))
	if errors.Is(err, entity.ErrNotFound) {
		return nil, entity.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	var tokens *Tokens
	err = u.tx.WithTx(ctx, func(ctx context.Context) error {
		if err := u.tokens.Revoke(ctx, rt.ID); err != nil {
			return err
		}
		t, err := u.issueTokens(ctx, rt.UserID, meta)
		if err != nil {
			return err
		}
		tokens = t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (u *Usecase) Logout(ctx context.Context, rawRefresh string) error {
	rt, err := u.tokens.GetActiveByHash(ctx, hashToken(rawRefresh))
	if errors.Is(err, entity.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return u.tokens.Revoke(ctx, rt.ID)
}

func (u *Usecase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.users.GetByID(ctx, userID)
}

func (u *Usecase) issueTokens(ctx context.Context, userID uuid.UUID, meta Meta) (*Tokens, error) {
	access, err := u.jwt.Generate(userID)
	if err != nil {
		return nil, err
	}

	raw, err := newRefreshToken()
	if err != nil {
		return nil, err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	err = u.tokens.Create(ctx, &entity.RefreshToken{
		ID:        id,
		UserID:    userID,
		TokenHash: hashToken(raw),
		UserAgent: meta.UserAgent,
		IP:        meta.IP,
		ExpiresAt: time.Now().Add(u.refreshTTL),
	})
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  access,
		RefreshToken: raw,
		ExpiresIn:    int64(u.accessTTL.Seconds()),
	}, nil
}

func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
