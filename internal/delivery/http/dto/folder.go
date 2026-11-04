package dto

import (
	"time"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type CreateFolderRequest struct {
	Name     string     `json:"name" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type RenameFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

type MoveFolderRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
}

type FolderResponse struct {
	ID        uuid.UUID  `json:"id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewFolderResponse(f *entity.Folder) FolderResponse {
	return FolderResponse{
		ID:        f.ID,
		ParentID:  f.ParentID,
		Name:      f.Name,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func NewFolderListResponse(items []entity.Folder) []FolderResponse {
	out := make([]FolderResponse, 0, len(items))
	for i := range items {
		out = append(out, NewFolderResponse(&items[i]))
	}
	return out
}
