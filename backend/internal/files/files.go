package files

import (
	"bytecrate/internal/database"
	"bytecrate/internal/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterFilesRoutes(r *gin.RouterGroup) {
	files := r.Group("/files")
	files.POST("/upload", UploadFile)
}

func UploadFile(c *gin.Context) {
	userID := c.GetUint("userID")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	// Create per-user directory
	userDir := filepath.Join("uploads", fmt.Sprint(userID))
	os.MkdirAll(userDir, 0755)

	// Generate filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
	filePath := filepath.Join(userDir, filename)

	// Save to disk
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Save file metadata in database
	fileRecord := models.File{
		UserID: userID,
		Name:   fileHeader.Filename,
		Path:   filePath,
		Size:   fileHeader.Size,
		Type:   fileHeader.Header.Get("Content-Type"),
	}

	database.DB.Create(&fileRecord)

	c.JSON(http.StatusOK, gin.H{
		"message": "uploaded",
		"file":    fileRecord,
	})
}