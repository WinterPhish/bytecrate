package models

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID                    uint      `json:"id"`
	Email                 string    `json:"email"`
	PasswordHash          string    `json:"-"` // never return
	Username              string    `json:"username"`
	StorageQuotaBytesUsed int64     `json:"storage_quota_bytes_used"`
	StorageQuotaBytes     int64     `json:"storage_quota_bytes"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
