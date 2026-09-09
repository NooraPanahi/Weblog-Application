package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxImageSize = 5 << 20 // 5 MB
	ImageDir     = "statics/images"
)

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

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create image file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("save image: %w", err)
	}

	return "/static/images/"+filename, nil
}
