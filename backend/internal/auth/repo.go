package auth

import (
	"bytecrate/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByID returns a user by numeric ID
func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) CreateUser(u *models.User) error {
	return r.DB.Create(u).Error
}

// SaveRefreshToken persists a refresh token record
func (r *UserRepo) SaveRefreshToken(t *models.RefreshToken) error {
	return r.DB.Create(t).Error
}

// FindRefreshToken looks up a refresh token by token string
func (r *UserRepo) FindRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.DB.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func (r *UserRepo) RevokeRefreshToken(id uint) error {
	return r.DB.Model(&models.RefreshToken{}).Where("id = ?", id).Updates(map[string]any{"revoked": true}).Error
}

// DeleteExpiredRefreshTokens removes tokens that are expired or revoked (optional cleanup)
func (r *UserRepo) DeleteExpiredRefreshTokens() error {
	return r.DB.Where("expires_at <= ? OR revoked = ?", time.Now(), true).Delete(&models.RefreshToken{}).Error
}
