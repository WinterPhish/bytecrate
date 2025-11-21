package http

import (
	docs "bytecrate/internal/http/docs"
	"bytecrate/internal/http/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	handlers.RegisterAuthRoutes(api)
	handlers.RegisterFilesRoutes(api)
	handlers.RegisterDevRoutes(api)

	docs.SwaggerInfo.BasePath = "/api/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
