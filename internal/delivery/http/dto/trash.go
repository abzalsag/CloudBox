package dto

import (
	"time"

	"CloudBox/internal/entity"
)

type TrashedFolderResponse struct {
	FolderResponse
	DeletedAt *time.Time `json:"deleted_at"`
}

type TrashedFileResponse struct {
	FileResponse
	DeletedAt *time.Time `json:"deleted_at"`
}

type TrashResponse struct {
	Folders []TrashedFolderResponse `json:"folders"`
	Files   []TrashedFileResponse   `json:"files"`
}

func NewTrashResponse(folders []entity.Folder, files []entity.File) TrashResponse {
	out := TrashResponse{
		Folders: make([]TrashedFolderResponse, 0, len(folders)),
		Files:   make([]TrashedFileResponse, 0, len(files)),
	}
	for i := range folders {
		out.Folders = append(out.Folders, TrashedFolderResponse{
			FolderResponse: NewFolderResponse(&folders[i]),
			DeletedAt:      folders[i].DeletedAt,
		})
	}
	for i := range files {
		out.Files = append(out.Files, TrashedFileResponse{
			FileResponse: NewFileResponse(&files[i]),
			DeletedAt:    files[i].DeletedAt,
		})
	}
	return out
}
