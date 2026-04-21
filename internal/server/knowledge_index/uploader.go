package knowledge_index

import (
	errorsx "OnCallAgent/pkg/errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileUploader interface {
	Upload(filename string, data []byte) (string, error)
}

type UploaderConfig struct {
	StoragePath string
	MaxSize     int64
}

type uploader struct {
	storagePath string
	maxSize     int64
}

func NewUploader(cfg UploaderConfig) FileUploader {
	return &uploader{
		storagePath: cfg.StoragePath,
		maxSize:     cfg.MaxSize,
	}
}

func (u *uploader) Upload(filename string, data []byte) (string, error) {
	if int64(len(data)) > u.maxSize {
		return "", errorsx.ErrFileTooLarge
	}

	absPath, err := filepath.Abs(u.storagePath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errorsx.ErrSaveFailed, err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return "", fmt.Errorf("%w: %v", errorsx.ErrSaveFailed, err)
	}

	dst := filepath.Join(absPath, filename)
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return "", fmt.Errorf("%w: %v", errorsx.ErrSaveFailed, err)
	}

	return dst, nil
}
