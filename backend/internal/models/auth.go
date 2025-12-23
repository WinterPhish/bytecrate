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

// RefreshToken model represents stored refresh tokens
type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	Token     string    `json:"token" gorm:"unique;size:512"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}
