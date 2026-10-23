package entity

import (
	"sort"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     sort.StringSlice
	StorageQuota int64
	StorageUser  int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
