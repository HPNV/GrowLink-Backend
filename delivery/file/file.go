package file

import (
	"net/http"

	"github.com/HPNV/growlink-backend/service"
	"github.com/gin-gonic/gin"
)

type IFile interface {
	UploadImage(c *gin.Context)
	GetByUUID(c *gin.Context)
	Delete(c *gin.Context)
	GetByUploadedBy(c *gin.Context)
}

type File struct {
	service service.IRegistry
}

func NewFile(service service.IRegistry) IFile {
	return &File{
		service: service,
	}
}

func (f *File) UploadImage(c *gin.Context) {
	// Get file from form
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Get uploaded_by from query param or authenticated user
	uploadedBy := c.Query("uploaded_by")
	if uploadedBy == "" {
		// Could get from auth middleware context
		userUUID, exists := c.Get("user_uuid")
		if exists {
			uploadedBy = userUUID.(string)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uploaded_by parameter required"})
			return
		}
	}

	// Validate file size (10MB limit)
	const maxFileSize = 10 << 20 // 10MB
	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum 10MB allowed"})
		return
	}

	result, err := f.service.GetFile().UploadImage(file, header, uploadedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (f *File) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	file, err := f.service.GetFile().GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, file)
}

func (f *File) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	err := f.service.GetFile().Delete(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (f *File) GetByUploadedBy(c *gin.Context) {
	uploadedBy := c.Param("uploadedBy")

	files, err := f.service.GetFile().GetByUploadedBy(uploadedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}
