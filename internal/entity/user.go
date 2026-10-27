package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	StorageQuota int64
	StorageUsed  int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
