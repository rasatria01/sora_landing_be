package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

const (
	MaxWidth  = 2048 // Maximum width for any image
	MaxHeight = 2048 // Maximum height for any image
	Quality   = 85   // JPEG quality (1-100)
)

// ProcessImage optimizes the image while maintaining quality
func ProcessImage(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read image into memory
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get current dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate new dimensions while maintaining aspect ratio
	if width > MaxWidth || height > MaxHeight {
		if width > height {
			height = height * MaxWidth / width
			width = MaxWidth
		} else {
			width = width * MaxHeight / height
			height = MaxHeight
		}
	}

	// Resize if needed
	if width != bounds.Dx() || height != bounds.Dy() {
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	// Create buffer for processed image
	var buf bytes.Buffer

	// Encode with appropriate format and quality
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: Quality})
	case "png":
		err = png.Encode(&buf, img)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// IsImage checks if the file is an image based on extension
func IsImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
}

// GetImageFormat returns the image format from filename
func GetImageFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	case ".webp":
		return "webp"
	default:
		return ""
	}
}

// SaveOptimizedImage saves the optimized image to a writer
func SaveOptimizedImage(src io.Reader, dst io.Writer, format string) error {
	// Decode image
	img, _, err := image.Decode(src)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Get dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Resize if needed
	if width > MaxWidth || height > MaxHeight {
		if width > height {
			height = height * MaxWidth / width
			width = MaxWidth
		} else {
			width = width * MaxHeight / height
			height = MaxHeight
		}
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	// Encode with appropriate format and quality
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(dst, img, &jpeg.Options{Quality: Quality})
	case "png":
		err = png.Encode(dst, img)
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}

	return err
}
