package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.POST("/register", Register)
	//auth.POST("/login", Login)
}

// Handler struct with repo
type AuthHandler struct {
	Repo *UserRepo
}

type UserRepo struct {
	DB *gorm.DB
}

// @Summary Register user
// @Description Register user to database
// @Accept  json
// @Produce  json
// @Router /auth/register [post]
func Register(c *gin.Context) {
	c.JSON(200, gin.H{"Register": "Reached"})
}
