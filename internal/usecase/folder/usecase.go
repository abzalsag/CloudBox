package folder

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Create(ctx context.Context, ownerID uuid.UUID, parentID *uuid.UUID, name string) (*entity.Folder, error) {
	name, err := validateName(name)
	if err != nil {
		return nil, err
	}
	if parentID != nil {
		if _, err := u.repo.GetByID(ctx, ownerID, *parentID); err != nil {
			return nil, err
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	f := &entity.Folder{ID: id, OwnerID: ownerID, ParentID: parentID, Name: name}
	if err := u.repo.Create(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (u *Usecase) Get(ctx context.Context, ownerID, id uuid.UUID) (*entity.Folder, error) {
	return u.repo.GetByID(ctx, ownerID, id)
}

func (u *Usecase) List(ctx context.Context, ownerID uuid.UUID, parentID *uuid.UUID, limit, offset int) ([]entity.Folder, error) {
	if parentID != nil {
		if _, err := u.repo.GetByID(ctx, ownerID, *parentID); err != nil {
			return nil, err
		}
	}
	return u.repo.ListByParent(ctx, ownerID, parentID, limit, offset)
}

func (u *Usecase) Rename(ctx context.Context, ownerID, id uuid.UUID, name string) (*entity.Folder, error) {
	name, err := validateName(name)
	if err != nil {
		return nil, err
	}
	return u.repo.Rename(ctx, ownerID, id, name)
}

func (u *Usecase) Move(ctx context.Context, ownerID, id uuid.UUID, newParentID *uuid.UUID) (*entity.Folder, error) {
	current, err := u.repo.GetByID(ctx, ownerID, id)
	if err != nil {
		return nil, err
	}

	if newParentID != nil {
		if *newParentID == id {
			return nil, fmt.Errorf("%w: cannot move folder into itself", entity.ErrInvalidInput)
		}
		if _, err := u.repo.GetByID(ctx, ownerID, *newParentID); err != nil {
			return nil, err
		}
		isDescendant, err := u.repo.IsDescendant(ctx, ownerID, *newParentID, id)
		if err != nil {
			return nil, err
		}
		if isDescendant {
			return nil, fmt.Errorf("%w: cannot move folder into its own subfolder", entity.ErrInvalidInput)
		}
	}

	if sameParent(current.ParentID, newParentID) {
		return current, nil
	}
	return u.repo.Move(ctx, ownerID, id, newParentID)
}

func (u *Usecase) Delete(ctx context.Context, ownerID, id uuid.UUID) error {
	if _, err := u.repo.GetByID(ctx, ownerID, id); err != nil {
		return err
	}
	return u.repo.SoftDeleteTree(ctx, ownerID, id)
}

func sameParent(a, b *uuid.UUID) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
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
