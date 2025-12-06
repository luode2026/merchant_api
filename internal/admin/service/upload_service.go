package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UploadService struct {
	ctx context.Context
}

func NewUploadService(ctx context.Context) *UploadService {
	return &UploadService{ctx: ctx}
}

// UploadImage handles the image upload logic
func (s *UploadService) UploadImage(file *multipart.FileHeader) (string, error) {
	// 1. Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return "", fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 2. Validate file size (e.g., max 5MB)
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("文件大小超过限制 (5MB)")
	}

	// 3. Generate unique filename
	filename := uuid.New().String() + ext

	// 4. Create upload directory if not exists
	// Organize by date to avoid too many files in one directory
	dateDir := time.Now().Format("20060102")
	uploadDir := filepath.Join("uploads", "images", dateDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 5. Save file
	dst := filepath.Join(uploadDir, filename)
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 6. Return relative path or URL
	// Assuming the static file server will serve from /uploads
	return "/" + dst, nil
}
