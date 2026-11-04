package dto

import (
	"github.com/google/uuid"

	"CloudBox/internal/entity"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	StorageQuota int64     `json:"storage_quota"`
	StorageUsed  int64     `json:"storage_used"`
}

func NewUserResponse(u *entity.User) UserResponse {
	return UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		FullName:     u.FullName,
		StorageQuota: u.StorageQuota,
		StorageUsed:  u.StorageUsed,
	}
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
