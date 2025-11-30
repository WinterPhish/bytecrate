package dev

import "github.com/gin-gonic/gin"

func RegisterDevRoutes(r *gin.RouterGroup) {
	files := r.Group("/dev")
	files.GET("/status", StatusCheck)
}

//
// @Summary Test for life
// @Description Dev test for backend
// @Accept  json
// @Produce  json
// @Router /dev/status [get]
func StatusCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
