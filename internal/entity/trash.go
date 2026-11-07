package entity

import "github.com/google/uuid"

type FileRef struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID
	Size       int64
	StorageKey string
}

type FolderRef struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
}
