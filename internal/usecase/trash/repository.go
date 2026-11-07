package trash

import (
	"context"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type Repository interface {
	ListTrashedFolders(ctx context.Context, ownerID uuid.UUID) ([]entity.Folder, error)
	ListTrashedFiles(ctx context.Context, ownerID uuid.UUID) ([]entity.File, error)
	GetFolderAny(ctx context.Context, ownerID, id uuid.UUID) (*entity.Folder, error)
	GetFileAny(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, error)
	RestoreFile(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, error)
	FolderTreeIDs(ctx context.Context, ownerID, id uuid.UUID) ([]uuid.UUID, error)
	RestoreFolderTree(ctx context.Context, ownerID uuid.UUID, ids []uuid.UUID) error
	FilesInFolders(ctx context.Context, ownerID uuid.UUID, ids []uuid.UUID) ([]entity.FileRef, error)
	HardDeleteFile(ctx context.Context, ownerID, id uuid.UUID) error
	HardDeleteFolders(ctx context.Context, ownerID uuid.UUID, ids []uuid.UUID) error
	ListExpiredFiles(ctx context.Context, limit int) ([]entity.FileRef, error)
	ListExpiredFolders(ctx context.Context, limit int) ([]entity.FolderRef, error)
}

type FolderReader interface {
	GetByID(ctx context.Context, ownerID, id uuid.UUID) (*entity.Folder, error)
}

type UserRepository interface {
	DecreaseStorageUsed(ctx context.Context, id uuid.UUID, delta int64) error
}

type Storage interface {
	Delete(ctx context.Context, key string) error
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
