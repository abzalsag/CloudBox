package folder

import (
	"context"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, f *entity.Folder) error
	GetByID(ctx context.Context, ownerID, id uuid.UUID) (*entity.Folder, error)
	ListByParent(ctx context.Context, ownerID uuid.UUID, parentID *uuid.UUID, limit, offset int) ([]entity.Folder, error)
	Rename(ctx context.Context, ownerID, id uuid.UUID, name string) (*entity.Folder, error)
	Move(ctx context.Context, ownerID, id uuid.UUID, parentID *uuid.UUID) (*entity.Folder, error)
	IsDescendant(ctx context.Context, ownerID, candidateID, ancestorID uuid.UUID) (bool, error)
	SoftDeleteTree(ctx context.Context, ownerID, id uuid.UUID) error
}
