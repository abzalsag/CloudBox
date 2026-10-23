package entity

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID
	FolderID   *uuid.UUID
	Name       string
	MimeType   string
	CheckSum   string
	StorageKey string
	Bucket     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
