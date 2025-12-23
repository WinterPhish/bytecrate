package auth

import (
	"bytecrate/internal/database"
	"bytecrate/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Handler struct with repo
type AuthHandler struct {
	Repo *UserRepo
}

func NewAuthHandler(repo *UserRepo) *AuthHandler {
	return &AuthHandler{Repo: repo}
}

func RegisterAuthRoutes(r *gin.RouterGroup) {
	userRepo := NewUserRepo(database.DB)
	authHandler := NewAuthHandler(userRepo)
	auth := r.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)
}

// Register creates a new user and returns a JWT
func (h *AuthHandler) Register(c *gin.Context) {
	var body models.RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check existence
	if _, err := h.Repo.FindByEmail(body.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &models.User{
		Email:             body.Email,
		Username:          body.Email,
		PasswordHash:      string(hashed),
		StorageQuotaBytes: 1073741824, // 1 GB
	}

	if err := h.Repo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// create refresh token and set cookie
	refreshStr, err := GenerateRefreshTokenString()
	if err == nil {
		expiresAt := time.Now().Add(14 * 24 * time.Hour)
		rt := &models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshStr,
			ExpiresAt: expiresAt,
		}
		_ = h.Repo.SaveRefreshToken(rt)
		secure := os.Getenv("COOKIE_SECURE") == "true"
		maxAge := int(time.Until(expiresAt).Seconds())
		c.SetCookie("refresh_token", refreshStr, maxAge, "/", "", secure, true)
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// Login validates credentials and returns JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var body models.LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repo.FindByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// create refresh token and set cookie
	refreshStr, err := GenerateRefreshTokenString()
	if err == nil {
		expiresAt := time.Now().Add(14 * 24 * time.Hour)
		rt := &models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshStr,
			ExpiresAt: expiresAt,
		}
		_ = h.Repo.SaveRefreshToken(rt)
		secure := os.Getenv("COOKIE_SECURE") == "true"
		maxAge := int(time.Until(expiresAt).Seconds())
		c.SetCookie("refresh_token", refreshStr, maxAge, "/", "", secure, true)
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// Refresh exchanges a valid refresh cookie for a new access token and rotates the refresh token
func (h *AuthHandler) Refresh(c *gin.Context) {
	rtCookie, err := c.Cookie("refresh_token")
	if err != nil || rtCookie == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	stored, err := h.Repo.FindRefreshToken(rtCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if stored.Revoked || stored.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired or revoked"})
		return
	}

	user, err := h.Repo.FindByID(stored.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// generate new access token
	access, err := GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	// rotate refresh token: create new one, revoke old
	newRefresh, err := GenerateRefreshTokenString()
	if err == nil {
		expiresAt := time.Now().Add(14 * 24 * time.Hour)
		newRT := &models.RefreshToken{
			UserID:    user.ID,
			Token:     newRefresh,
			ExpiresAt: expiresAt,
		}
		_ = h.Repo.SaveRefreshToken(newRT)
		_ = h.Repo.RevokeRefreshToken(stored.ID)
		secure := os.Getenv("COOKIE_SECURE") == "true"
		maxAge := int(time.Until(expiresAt).Seconds())
		c.SetCookie("refresh_token", newRefresh, maxAge, "/", "", secure, true)
	}

	c.JSON(http.StatusOK, gin.H{"token": access})
}
