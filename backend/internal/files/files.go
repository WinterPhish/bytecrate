package files

import (
	"bytecrate/internal/database"
	"bytecrate/internal/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterFilesRoutes(r *gin.RouterGroup) {
	files := r.Group("/files")
	files.POST("/upload", UploadFile)
	files.GET("/list/", ListFiles)
	files.GET("/download/:id", DownloadFile)
}

func UploadFile(c *gin.Context) {
	userID := c.GetUint("userID")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

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

	database.DB.Create(&fileRecord)

	c.JSON(http.StatusOK, gin.H{
		"message": "uploaded",
		"file":    fileRecord,
	})
}

// TODO : MAKE USERID PARAM OPTIONAL AND ONLY ALLOW IF ADMIN

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
