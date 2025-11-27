package internal

import (
	"bytecrate/internal/auth"
	"bytecrate/internal/dev"
	docs "bytecrate/internal/docs"
	"bytecrate/internal/files"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default()) // All origins allowed by default

	api := r.Group("/api")
	auth.RegisterAuthRoutes(api, db)
	files.RegisterFilesRoutes(api)
	dev.RegisterDevRoutes(api)

	docs.SwaggerInfo.BasePath = "/api/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
