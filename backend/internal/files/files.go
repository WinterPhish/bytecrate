package files

import (
	"bytecrate/internal/database"
	"bytecrate/internal/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterFilesRoutes(r *gin.RouterGroup) {
	files := r.Group("/files")
	files.POST("/upload", UploadFile)
	files.GET("/list/", ListFiles)
	files.DELETE("/:id", DeleteFile)
	files.PUT("/:id/rename", RenameFile)
	files.GET("/:id/download", DownloadFile)
}

func UploadFile(c *gin.Context) {
	userID := c.GetUint("userID")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	var user models.User

	if err := database.DB.Where("id = ?", userID).Order("created_at DESC").Find(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	if fileHeader.Size > user.StorageQuotaBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds quota limit"})
		return
	}

	if user.StorageQuotaBytesUsed+fileHeader.Size > user.StorageQuotaBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "upload would exceed your storage quota"})
		return
	}

	fmt.Println("Storage used: ", user.StorageQuotaBytesUsed, " / ", user.StorageQuotaBytes)

	// Create per-user directory
	saveDir := "/app/uploads/" + strconv.Itoa(int(userID))
	os.MkdirAll(saveDir, 0755)

	// Generate filename
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
	filePath := filepath.Join(saveDir, fileName)

	// Save to disk
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	fmt.Println("Saved file to ", filePath)

	// Save file metadata in database
	fileRecord := models.File{
		UserID:      userID,
		Filename:    fileHeader.Filename,
		Path:        filePath,
		SizeBytes:   fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}

	if err := database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("storage_quota_bytes_used", gorm.Expr("storage_quota_bytes_used + ?", fileHeader.Size)).
		Error; err != nil {
		return
	}

	database.DB.Create(&fileRecord)

	c.JSON(http.StatusOK, gin.H{
		"message": "uploaded",
		"file":    fileRecord,
	})
}

func DeleteFile(c *gin.Context) {
	db := database.DB
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Delete file from disk if exists
	if _, err := os.Stat(file.Path); err == nil {
		_ = os.Remove(file.Path)
	}

	// Remove DB record
	if err := db.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from DB"})
		return
	}

	// Subtract quota
	db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("storage_quota_bytes_used",
			gorm.Expr("storage_quota_bytes_used - ?", file.SizeBytes),
		)

	c.JSON(http.StatusOK, gin.H{"message": "File deleted"})
}

type RenameRequest struct {
	Filename string `json:"filename"`
}

func RenameFile(c *gin.Context) {
	db := database.DB
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var body RenameRequest
	if err := c.BindJSON(&body); err != nil || strings.TrimSpace(body.Filename) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	file.Filename = body.Filename

	if err := db.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rename"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Renamed", "file": file})
}

func ListFiles(c *gin.Context) {
	userID := c.GetUint("userID")
	fmt.Println("Listing files for user: ", userID)

	var files []models.File
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func DownloadFile(c *gin.Context) {
	db := database.DB
	userID := c.GetUint("userID")
	fileID := c.Param("id")
	fmt.Println("User ", userID, " downloading file ", fileID)

	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		abs, _ := filepath.Abs(file.Path)
		fmt.Println("DEBUG: Cannot find file. DB path =", file.Path, "Absolute =", abs)
		c.JSON(http.StatusNotFound, gin.H{"error": "File missing from server"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Filename))
	c.Header("Content-Type", file.ContentType)
	c.File(file.Path)
}
