package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxImageSize = 5 << 20 // 5 MB
	ImageDir     = "statics/images"
)

var allowedContentTypes = map[string]bool{
    "image/jpeg": true,
    "image/png":  true,
    "image/gif":  true,
    "image/webp": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

func SaveImage(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", nil
	}
	if file.Size > MaxImageSize {
		return "", fmt.Errorf("image too large")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedExtensions[ext] {
		return "", fmt.Errorf("invalid image type")
	}

	if err := os.MkdirAll(ImageDir, 0755); err != nil {
		return "", fmt.Errorf("create image directory: %w", err)
	}

	filename := uuid.New().String() + ext
	filePath := filepath.Join(ImageDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open uploaded image: %w", err)
	}
	defer src.Close()

	buffer := make([]byte, 512)

	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read uploaded image: %w", err)
	}

	contentType := http.DetectContentType(buffer)

	if !allowedContentTypes[contentType] {
		return "", fmt.Errorf("invalid image content type")
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create image file: %w", err)
	}
	defer dst.Close()

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("reset uploaded image: %w", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("save image: %w", err)
	}

	return "/static/images/" + filename, nil
}

func DeleteImage(imagePath string) error {
	if imagePath == "" {
		return nil
	}

	const urlPrefix = "/static/images/"

	if !strings.HasPrefix(imagePath, urlPrefix) {
		return fmt.Errorf("invalid image path")
	}

	filename := strings.TrimPrefix(imagePath, urlPrefix)

	filePath := filepath.Join(ImageDir, filename)

	if err := os.Remove(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("delete image: %w", err)
	}

	return nil
}
