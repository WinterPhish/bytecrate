package dev

import "github.com/gin-gonic/gin"

func RegisterDevRoutes(r *gin.RouterGroup) {
	dev := r.Group("/dev")
	dev.GET("/status", StatusCheck)
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
