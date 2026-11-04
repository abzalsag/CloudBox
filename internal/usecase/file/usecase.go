package file

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

const defaultContentType = "application/octet-stream"

type Usecase struct {
	repo    Repository
	folders FolderReader
	users   UserRepository
	storage Storage
	tx      TxManager
}

func NewUsecase(repo Repository, folders FolderReader, users UserRepository, storage Storage, tx TxManager) *Usecase {
	return &Usecase{repo: repo, folders: folders, users: users, storage: storage, tx: tx}
}

type UploadInput struct {
	OwnerID     uuid.UUID
	FolderID    *uuid.UUID
	Name        string
	ContentType string
	Content     io.ReadSeeker
}

func (u *Usecase) Upload(ctx context.Context, in UploadInput) (*entity.File, error) {
	name, err := validateName(in.Name)
	if err != nil {
		return nil, err
	}
	if in.FolderID != nil {
		if _, err := u.folders.GetByID(ctx, in.OwnerID, *in.FolderID); err != nil {
			return nil, err
		}
	}

	checksum, size, err := checksumAndSize(in.Content)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetByID(ctx, in.OwnerID)
	if err != nil {
		return nil, err
	}
	if user.StorageUsed+size > user.StorageQuota {
		return nil, entity.ErrQuotaExceeded
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s/%s", in.OwnerID, id)

	contentType := in.ContentType
	if contentType == "" {
		contentType = defaultContentType
	}

	if err := u.storage.Put(ctx, key, in.Content, size, contentType); err != nil {
		return nil, err
	}

	f := &entity.File{
		ID:         id,
		OwnerID:    in.OwnerID,
		FolderID:   in.FolderID,
		Name:       name,
		Size:       size,
		MimeType:   contentType,
		Checksum:   checksum,
		StorageKey: key,
		Bucket:     u.storage.Bucket(),
	}

	err = u.tx.WithTx(ctx, func(ctx context.Context) error {
		if err := u.repo.Create(ctx, f); err != nil {
			return err
		}
		return u.users.IncreaseStorageUsed(ctx, in.OwnerID, size)
	})
	if err != nil {
		_ = u.storage.Delete(context.WithoutCancel(ctx), key)
		return nil, err
	}
	return f, nil
}

func (u *Usecase) Get(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, error) {
	return u.repo.GetByID(ctx, ownerID, id)
}

func (u *Usecase) List(ctx context.Context, ownerID uuid.UUID, folderID *uuid.UUID, limit, offset int) ([]entity.File, error) {
	if folderID != nil {
		if _, err := u.folders.GetByID(ctx, ownerID, *folderID); err != nil {
			return nil, err
		}
	}
	return u.repo.ListByFolder(ctx, ownerID, folderID, limit, offset)
}

func (u *Usecase) Download(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, io.ReadCloser, error) {
	f, err := u.repo.GetByID(ctx, ownerID, id)
	if err != nil {
		return nil, nil, err
	}
	body, err := u.storage.Get(ctx, f.StorageKey)
	if err != nil {
		return nil, nil, err
	}
	return f, body, nil
}

func (u *Usecase) Rename(ctx context.Context, ownerID, id uuid.UUID, name string) (*entity.File, error) {
	name, err := validateName(name)
	if err != nil {
		return nil, err
	}
	return u.repo.Rename(ctx, ownerID, id, name)
}

func (u *Usecase) Move(ctx context.Context, ownerID, id uuid.UUID, folderID *uuid.UUID) (*entity.File, error) {
	if folderID != nil {
		if _, err := u.folders.GetByID(ctx, ownerID, *folderID); err != nil {
			return nil, err
		}
	}
	return u.repo.Move(ctx, ownerID, id, folderID)
}

func (u *Usecase) Delete(ctx context.Context, ownerID, id uuid.UUID) error {
	return u.repo.SoftDelete(ctx, ownerID, id)
}

func checksumAndSize(r io.ReadSeeker) (string, int64, error) {
	h := sha256.New()
	n, err := io.Copy(h, r)
	if err != nil {
		return "", 0, err
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}

func validateName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	switch {
	case name == "":
		return "", fmt.Errorf("%w: name is required", entity.ErrInvalidInput)
	case utf8.RuneCountInString(name) > 255:
		return "", fmt.Errorf("%w: name is too long", entity.ErrInvalidInput)
	case strings.ContainsAny(name, `/\`):
		return "", fmt.Errorf("%w: name must not contain slashes", entity.ErrInvalidInput)
	case name == "." || name == "..":
		return "", fmt.Errorf("%w: invalid name", entity.ErrInvalidInput)
	}
	return name, nil
}
