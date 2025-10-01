package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sora_landing_be/pkg/errors"
	"sora_landing_be/pkg/http/server/http_response"
	"sora_landing_be/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileController struct {
}

func NewFileController() FileController {
	return FileController{}
}

// UploadFiles handles multiple file uploads
func (ctl *FileController) UploadFile(c *gin.Context) {
	// Limit file size (in MB → bytes)
	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		3024*1024,
	)

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		http_response.SendError(c, errors.StorageErrorToAppError("Failed to read uploaded file"))
		return
	}

	// Validate file format
	if ok, expectedFormat, actualFormat := utils.IsDocumentFile(fileHeader.Filename); !ok {
		http_response.SendError(c, errors.StorageErrorToAppError(
			fmt.Sprintf("Cannot upload file with format: %s, expected: %s", actualFormat, expectedFormat),
		))
		return
	}

	// Generate a safe/unique filename
	filename := utils.GenerateKeyFile(fileHeader.Filename)

	// Save file to local storage (e.g., ./uploads)
	savePath := path.Join("uploads", filename)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		http_response.SendError(c, errors.StorageErrorToAppError("Failed to save file"))
		return
	}

	// Return success response
	http_response.SendSuccess(c, http.StatusOK, "Success upload file", map[string]interface{}{
		"filename": filename,
		"path":     savePath,
	})
}

// GetPublicFile gets a public file by ID
func (ctl *FileController) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		http_response.SendError(c, errors.StorageErrorToAppError("Filename is required"))
		return
	}

	// Clean the filename to prevent ../../ attacks
	cleanName := filepath.Base(filename) // keeps only the last element
	uploadDir := "uploads"
	filePath := filepath.Join(uploadDir, cleanName)

	// Extra safety: ensure the resolved path is still inside uploadDir
	absUploadDir, _ := filepath.Abs(uploadDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absUploadDir) {
		http_response.SendError(c, errors.StorageErrorToAppError("Invalid file path"))
		return
	}

	// Try to remove file
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			http_response.SendError(c, errors.StorageErrorToAppError("File not found"))
			return
		}
		http_response.SendError(c, errors.StorageErrorToAppError("Failed to delete file"))
		return
	}

	http_response.SendSuccess(c, http.StatusOK, "File deleted successfully", nil)
}

// DeleteFile deletes a file by ID
