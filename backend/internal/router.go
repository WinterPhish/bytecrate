package internal

import (
	"bytecrate/internal/auth"
	"bytecrate/internal/dev"
	docs "bytecrate/internal/docs"
	"bytecrate/internal/files"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	auth.RegisterAuthRoutes(api)
	files.RegisterFilesRoutes(api)
	dev.RegisterDevRoutes(api)

	docs.SwaggerInfo.BasePath = "/api/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
