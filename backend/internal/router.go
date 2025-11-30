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
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length", "Authorization"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))
	api := r.Group("/api")
	auth.RegisterAuthRoutes(api, db)

	// Protected routes with auth middleware
	protected := api.Group("")
	protected.Use(auth.AuthMiddleware())
	files.RegisterFilesRoutes(protected)
	dev.RegisterDevRoutes(protected)

	docs.SwaggerInfo.BasePath = "/api/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
