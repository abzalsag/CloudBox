package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"CloudBox/internal/entity"
	"CloudBox/internal/repository/postgres/sqlcgen"
)

type FileRepo struct {
	*Store
}

func NewFileRepo(s *Store) *FileRepo { return &FileRepo{Store: s} }

func (r *FileRepo) Create(ctx context.Context, f *entity.File) error {
	row, err := r.Q(ctx).CreateFile(ctx, sqlcgen.CreateFileParams{
		ID:         f.ID,
		OwnerID:    f.OwnerID,
		FolderID:   toNullUUID(f.FolderID),
		Name:       f.Name,
		Size:       f.Size,
		MimeType:   f.MimeType,
		Checksum:   f.Checksum,
		StorageKey: f.StorageKey,
		Bucket:     f.Bucket,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return entity.ErrAlreadyExists
		}
		return err
	}
	*f = *toFile(row)
	return nil
}

func (r *FileRepo) GetByID(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, error) {
	row, err := r.Q(ctx).GetFile(ctx, sqlcgen.GetFileParams{ID: id, OwnerID: ownerID})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toFile(row), nil
}

func (r *FileRepo) ListByFolder(ctx context.Context, ownerID uuid.UUID, folderID *uuid.UUID, limit, offset int) ([]entity.File, error) {
	rows, err := r.Q(ctx).ListFilesByFolder(ctx, sqlcgen.ListFilesByFolderParams{
		OwnerID:  ownerID,
		FolderID: toNullUUID(folderID),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}
	out := make([]entity.File, 0, len(rows))
	for _, row := range rows {
		out = append(out, *toFile(row))
	}
	return out, nil
}

func (r *FileRepo) Rename(ctx context.Context, ownerID, id uuid.UUID, name string) (*entity.File, error) {
	row, err := r.Q(ctx).RenameFile(ctx, sqlcgen.RenameFileParams{ID: id, OwnerID: ownerID, Name: name})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return nil, entity.ErrAlreadyExists
		}
		return nil, err
	}
	return toFile(row), nil
}

func (r *FileRepo) Move(ctx context.Context, ownerID, id uuid.UUID, folderID *uuid.UUID) (*entity.File, error) {
	row, err := r.Q(ctx).MoveFile(ctx, sqlcgen.MoveFileParams{
		ID:       id,
		OwnerID:  ownerID,
		FolderID: toNullUUID(folderID),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return nil, entity.ErrAlreadyExists
		}
		return nil, err
	}
	return toFile(row), nil
}

func (r *FileRepo) SoftDelete(ctx context.Context, ownerID, id uuid.UUID) error {
	_, err := r.Q(ctx).SoftDeleteFile(ctx, sqlcgen.SoftDeleteFileParams{ID: id, OwnerID: ownerID})
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.ErrNotFound
	}
	return err
}

func toFile(row sqlcgen.File) *entity.File {
	return &entity.File{
		ID:         row.ID,
		OwnerID:    row.OwnerID,
		FolderID:   fromNullUUID(row.FolderID),
		Name:       row.Name,
		Size:       row.Size,
		MimeType:   row.MimeType,
		Checksum:   row.Checksum,
		StorageKey: row.StorageKey,
		Bucket:     row.Bucket,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
		DeletedAt:  row.DeletedAt,
	}
}
