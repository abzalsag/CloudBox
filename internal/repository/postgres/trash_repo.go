package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"CloudBox/internal/entity"
	"CloudBox/internal/repository/postgres/sqlcgen"
)

type TrashRepo struct {
	*Store
}

func NewTrashRepo(s *Store) *TrashRepo { return &TrashRepo{Store: s} }

func (r *TrashRepo) ListTrashedFolders(ctx context.Context, ownerID uuid.UUID) ([]entity.Folder, error) {
	rows, err := r.Q(ctx).ListTrashedRootFolders(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Folder, 0, len(rows))
	for _, row := range rows {
		out = append(out, *toFolder(row))
	}
	return out, nil
}

func (r *TrashRepo) ListTrashedFiles(ctx context.Context, ownerID uuid.UUID) ([]entity.File, error) {
	rows, err := r.Q(ctx).ListTrashedFiles(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]entity.File, 0, len(rows))
	for _, row := range rows {
		out = append(out, *toFile(row))
	}
	return out, nil
}

func (r *TrashRepo) GetFolderAny(ctx context.Context, ownerID, id uuid.UUID) (*entity.Folder, error) {
	row, err := r.Q(ctx).GetFolderIncludingDeleted(ctx, sqlcgen.GetFolderIncludingDeletedParams{ID: id, OwnerID: ownerID})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toFolder(row), nil
}

func (r *TrashRepo) GetFileAny(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, error) {
	row, err := r.Q(ctx).GetFileIncludingDeleted(ctx, sqlcgen.GetFileIncludingDeletedParams{ID: id, OwnerID: ownerID})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toFile(row), nil
}

func (r *TrashRepo) RestoreFile(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, error) {
	row, err := r.Q(ctx).RestoreFile(ctx, sqlcgen.RestoreFileParams{ID: id, OwnerID: ownerID})
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

func (r *TrashRepo) FolderTreeIDs(ctx context.Context, ownerID, id uuid.UUID) ([]uuid.UUID, error) {
	return r.Q(ctx).ListFolderTreeIDs(ctx, sqlcgen.ListFolderTreeIDsParams{ID: id, OwnerID: ownerID})
}

func (r *TrashRepo) RestoreFolderTree(ctx context.Context, ownerID uuid.UUID, ids []uuid.UUID) error {
	err := r.WithTx(ctx, func(ctx context.Context) error {
		q := r.Q(ctx)
		if err := q.RestoreFolders(ctx, sqlcgen.RestoreFoldersParams{OwnerID: ownerID, Ids: ids}); err != nil {
			return err
		}
		return q.RestoreFilesInFolders(ctx, sqlcgen.RestoreFilesInFoldersParams{OwnerID: ownerID, FolderIds: ids})
	})
	if err != nil && isUniqueViolation(err) {
		return entity.ErrAlreadyExists
	}
	return err
}

func (r *TrashRepo) FilesInFolders(ctx context.Context, ownerID uuid.UUID, ids []uuid.UUID) ([]entity.FileRef, error) {
	rows, err := r.Q(ctx).ListFilesInFolders(ctx, sqlcgen.ListFilesInFoldersParams{OwnerID: ownerID, FolderIds: ids})
	if err != nil {
		return nil, err
	}
	out := make([]entity.FileRef, 0, len(rows))
	for _, row := range rows {
		out = append(out, entity.FileRef{ID: row.ID, OwnerID: row.OwnerID, Size: row.Size, StorageKey: row.StorageKey})
	}
	return out, nil
}

func (r *TrashRepo) HardDeleteFile(ctx context.Context, ownerID, id uuid.UUID) error {
	return r.Q(ctx).HardDeleteFile(ctx, sqlcgen.HardDeleteFileParams{ID: id, OwnerID: ownerID})
}

func (r *TrashRepo) HardDeleteFolders(ctx context.Context, ownerID uuid.UUID, ids []uuid.UUID) error {
	return r.Q(ctx).HardDeleteFolders(ctx, sqlcgen.HardDeleteFoldersParams{OwnerID: ownerID, Ids: ids})
}

func (r *TrashRepo) ListExpiredFiles(ctx context.Context, limit int) ([]entity.FileRef, error) {
	rows, err := r.Q(ctx).ListExpiredTrashedFiles(ctx, int32(limit))
	if err != nil {
		return nil, err
	}
	out := make([]entity.FileRef, 0, len(rows))
	for _, row := range rows {
		out = append(out, entity.FileRef{ID: row.ID, OwnerID: row.OwnerID, Size: row.Size, StorageKey: row.StorageKey})
	}
	return out, nil
}

func (r *TrashRepo) ListExpiredFolders(ctx context.Context, limit int) ([]entity.FolderRef, error) {
	rows, err := r.Q(ctx).ListExpiredTrashedFolderIDs(ctx, int32(limit))
	if err != nil {
		return nil, err
	}
	out := make([]entity.FolderRef, 0, len(rows))
	for _, row := range rows {
		out = append(out, entity.FolderRef{ID: row.ID, OwnerID: row.OwnerID})
	}
	return out, nil
}
