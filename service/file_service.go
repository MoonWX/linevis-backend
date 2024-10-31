// service/file_service.go
package service

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileService struct {
	uploadDir string
}

// UploadResult 定义上传结果的结构
type UploadResult struct {
	FullPath string // 完整的文件路径
	FileName string // 新的文件名（包含扩展名）
	Size     int64  // 文件大小
}

type FileNameStrategy int

const (
	OriginalName FileNameStrategy = iota
	TimeStampPrefix
	UUIDPrefix
	CustomName
)

func NewFileService(uploadDir string) (*FileService, error) {
	fs := &FileService{
		uploadDir: uploadDir,
	}

	if err := fs.EnsureUploadDir(); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	return fs, nil
}

func (s *FileService) EnsureUploadDir() error {
	absPath, err := filepath.Abs(s.uploadDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	s.uploadDir = absPath
	return nil
}

func (s *FileService) GenerateFileName(originalName string, strategy FileNameStrategy, customName string) string {
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)

	switch strategy {
	case TimeStampPrefix:
		timestamp := time.Now().Format("20060102150405")
		return fmt.Sprintf("%s_%s%s", timestamp, baseName, ext)

	case UUIDPrefix:
		return fmt.Sprintf("%s%s", uuid.New().String(), ext)

	case CustomName:
		if customName != "" {
			if filepath.Ext(customName) == "" {
				return customName + ext
			}
			return customName
		}
		fallthrough

	default:
		return originalName
	}
}

func (s *FileService) SaveUploadedFileWithName(file *multipart.FileHeader, strategy FileNameStrategy, customName string) (*UploadResult, error) {
	if err := s.EnsureUploadDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure upload directory exists: %v", err)
	}

	newFilename := s.GenerateFileName(file.Filename, strategy, customName)
	dst := filepath.Join(s.uploadDir, newFilename)

	// 先创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		os.Remove(dst) // 清理已创建的文件
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// 使用缓冲区复制文件内容
	buffer := make([]byte, 32*1024) // 32KB 缓冲区
	written, err := io.CopyBuffer(dstFile, src, buffer)
	if err != nil {
		os.Remove(dst) // 清理已创建的文件
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	// 确保所有数据都已写入磁盘
	if err := dstFile.Sync(); err != nil {
		os.Remove(dst)
		return nil, fmt.Errorf("failed to sync file to disk: %v", err)
	}

	// 返回上传结果
	return &UploadResult{
		FullPath: dst,
		FileName: newFilename,
		Size:     written,
	}, nil
}

func (s *FileService) SaveUploadedFile(file *multipart.FileHeader) (*UploadResult, error) {
	return s.SaveUploadedFileWithName(file, OriginalName, "")
}
