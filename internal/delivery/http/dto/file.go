package dto

import (
	"time"

	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type RenameFileRequest struct {
	Name string `json:"name" binding:"required"`
}

type MoveFileRequest struct {
	FolderID *uuid.UUID `json:"folder_id"`
}

type FileResponse struct {
	ID        uuid.UUID  `json:"id"`
	FolderID  *uuid.UUID `json:"folder_id"`
	Name      string     `json:"name"`
	Size      int64      `json:"size"`
	MimeType  string     `json:"mime_type"`
	Checksum  string     `json:"checksum"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewFileResponse(f *entity.File) FileResponse {
	return FileResponse{
		ID:        f.ID,
		FolderID:  f.FolderID,
		Name:      f.Name,
		Size:      f.Size,
		MimeType:  f.MimeType,
		Checksum:  f.Checksum,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func NewFileListResponse(items []entity.File) []FileResponse {
	out := make([]FileResponse, 0, len(items))
	for i := range items {
		out = append(out, NewFileResponse(&items[i]))
	}
	return out
}
