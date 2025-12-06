package controller

import (
	"merchant_api/internal/admin/service"
	"merchant_api/internal/pkg/response"
	"merchant_api/pkg/config"

	"github.com/gin-gonic/gin"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadImage handles single image upload
func (ctrl *UploadController) UploadImage(c *gin.Context) {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequestWithKey(c, "error.upload.file_retrieval_failed", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	// Call service
	uploadService := service.NewUploadService(c.Request.Context())
	path, err := uploadService.UploadImage(file)
	if err != nil {
		response.BadRequestWithKey(c, "error.upload.failed", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	// Construct full URL
	domain := config.GlobalConfig.Server.Admin.Domain
	fullURL := domain + path

	// Return success response with file URL and path
	response.Success(c, gin.H{
		"path": path,
		"url":  fullURL,
	})
}
