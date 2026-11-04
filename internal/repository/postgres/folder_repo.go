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

type FolderRepo struct {
	*Store
}

func NewFolderRepo(s *Store) *FolderRepo { return &FolderRepo{Store: s} }

func (r *FolderRepo) Create(ctx context.Context, f *entity.Folder) error {
	row, err := r.Q(ctx).CreateFolder(ctx, sqlcgen.CreateFolderParams{
		ID:       f.ID,
		OwnerID:  f.OwnerID,
		ParentID: toNullUUID(f.ParentID),
		Name:     f.Name,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return entity.ErrAlreadyExists
		}
		return err
	}
	*f = *toFolder(row)
	return nil
}

func (r *FolderRepo) GetByID(ctx context.Context, ownerID, id uuid.UUID) (*entity.Folder, error) {
	row, err := r.Q(ctx).GetFolder(ctx, sqlcgen.GetFolderParams{ID: id, OwnerID: ownerID})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toFolder(row), nil
}

func (r *FolderRepo) ListByParent(ctx context.Context, ownerID uuid.UUID, parentID *uuid.UUID, limit, offset int) ([]entity.Folder, error) {
	rows, err := r.Q(ctx).ListFoldersByParent(ctx, sqlcgen.ListFoldersByParentParams{
		OwnerID:  ownerID,
		ParentID: toNullUUID(parentID),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}
	out := make([]entity.Folder, 0, len(rows))
	for _, row := range rows {
		out = append(out, *toFolder(row))
	}
	return out, nil
}

func (r *FolderRepo) Rename(ctx context.Context, ownerID, id uuid.UUID, name string) (*entity.Folder, error) {
	row, err := r.Q(ctx).RenameFolder(ctx, sqlcgen.RenameFolderParams{ID: id, OwnerID: ownerID, Name: name})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return nil, entity.ErrAlreadyExists
		}
		return nil, err
	}
	return toFolder(row), nil
}

func (r *FolderRepo) Move(ctx context.Context, ownerID, id uuid.UUID, parentID *uuid.UUID) (*entity.Folder, error) {
	row, err := r.Q(ctx).MoveFolder(ctx, sqlcgen.MoveFolderParams{
		ID:       id,
		OwnerID:  ownerID,
		ParentID: toNullUUID(parentID),
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
	return toFolder(row), nil
}

func (r *FolderRepo) IsDescendant(ctx context.Context, candidateID, ancestorID uuid.UUID) (bool, error) {
	return r.Q(ctx).IsFolderDescendant(ctx, sqlcgen.IsFolderDescendantParams{
		CandidateID: candidateID,
		AncestorID:  ancestorID,
	})
}

func (r *FolderRepo) SoftDeleteTree(ctx context.Context, ownerID, id uuid.UUID) error {
	return r.WithTx(ctx, func(ctx context.Context) error {
		q := r.Q(ctx)

		ids, err := q.ListFolderTreeIDs(ctx, sqlcgen.ListFolderTreeIDsParams{ID: id, OwnerID: ownerID})
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			return entity.ErrNotFound
		}

		if err := q.SoftDeleteFolders(ctx, sqlcgen.SoftDeleteFoldersParams{OwnerID: ownerID, Ids: ids}); err != nil {
			return err
		}
		return q.SoftDeleteFilesInFolders(ctx, sqlcgen.SoftDeleteFilesInFoldersParams{OwnerID: ownerID, FolderIds: ids})
	})
}

func toFolder(row sqlcgen.Folder) *entity.Folder {
	return &entity.Folder{
		ID:        row.ID,
		OwnerID:   row.OwnerID,
		ParentID:  fromNullUUID(row.ParentID),
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: row.DeletedAt,
	}
}

func toNullUUID(id *uuid.UUID) uuid.NullUUID {
	if id == nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{UUID: *id, Valid: true}
}

func fromNullUUID(n uuid.NullUUID) *uuid.UUID {
	if !n.Valid {
		return nil
	}
	id := n.UUID
	return &id
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
