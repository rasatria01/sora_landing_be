package controllers

import (
	"net/http"
	"sora_landing_be/cmd/services"
	"sora_landing_be/pkg/errors"
	internalHTTP "sora_landing_be/pkg/http"
	"sora_landing_be/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService services.FileService
}

func NewFileController(fileService services.FileService) FileController {
	return FileController{
		fileService: fileService,
	}
}

// UploadFiles handles multiple file uploads
func (ctl *FileController) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		http_response.SendError(c, errors.NewValidationError("Invalid form data", map[string]string{"form": err.Error()}))
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		http_response.SendError(c, errors.NewValidationError("No files uploaded", map[string]string{"files": "required"}))
		return
	}

	module := c.PostForm("module")
	if module == "" {
		module = "general" // Default module if not specified
	}

	// Get user ID from context (set by auth middleware if needed)
	// Default for public endpoint

	responses, err := ctl.fileService.UploadFiles(c, files, module)
	if err != nil {
		http_response.SendError(c, err)
		return
	}

	http_response.SendSuccess(c, http.StatusCreated, "Files uploaded successfully", responses)
}

// GetPublicFile gets a public file by ID
func (ctl *FileController) GetPublicFile(c *gin.Context) {
	id, err := internalHTTP.BindParams[string](c, "id")
	if err != nil {
		http_response.SendError(c, errors.NewValidationError("Invalid ID", map[string]string{"id": err.Error()}))
		return
	}

	file, err := ctl.fileService.GetPublicFile(c, id)
	if err != nil {
		http_response.SendError(c, err)
		return
	}

	c.File(file.FilePath)
}

// DeleteFile deletes a file by ID
func (ctl *FileController) DeleteFile(c *gin.Context) {
	id, err := internalHTTP.BindParams[string](c, "id")
	if err != nil {
		http_response.SendError(c, errors.NewValidationError("Invalid ID", map[string]string{"id": err.Error()}))
		return
	}

	err = ctl.fileService.DeleteFile(c, id)
	if err != nil {
		http_response.SendError(c, err)
		return
	}

	http_response.SendSuccess(c, http.StatusOK, "File deleted successfully", nil)
}
