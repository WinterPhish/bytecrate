package http

import (
    "github.com/gin-gonic/gin"
    "bytecrate/internal/http/handlers"
)

func NewRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/health", handlers.HealthCheck)

    return r
}
