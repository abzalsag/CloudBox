package file

import (
	"context"
	"io"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, f *entity.File) error
	GetByID(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, error)
	ListByFolder(ctx context.Context, ownerID uuid.UUID, folderID *uuid.UUID, limit, offset int) ([]entity.File, error)
	Rename(ctx context.Context, ownerID, id uuid.UUID, name string) (*entity.File, error)
	Move(ctx context.Context, ownerID, id uuid.UUID, folderID *uuid.UUID) (*entity.File, error)
	SoftDelete(ctx context.Context, ownerID, id uuid.UUID) error
}

type FolderReader interface {
	GetByID(ctx context.Context, ownerID, id uuid.UUID) (*entity.Folder, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	IncreaseStorageUsed(ctx context.Context, id uuid.UUID, delta int64) error
}

type Storage interface {
	Put(ctx context.Context, key string, r io.ReadSeeker, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	Bucket() string
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
