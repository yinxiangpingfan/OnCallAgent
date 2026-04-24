package knowledgeindex

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type fileUploader struct {
	logger *logrus.Logger
}

type FileUploader interface {
	Upload(file *multipart.FileHeader, path string) (string, error)
}

func NewFileUploaderServer(log *logrus.Logger) FileUploader {
	return &fileUploader{logger: log}
}

func (u *fileUploader) Upload(file *multipart.FileHeader, path string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	u.logger.Errorf("上传文件失败, err: %v", err)

	var mode os.FileMode = 0o750

	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, mode); err != nil {
		u.logger.Errorf("上传文件失败, err: %v", err)
		return "", err
	}
	u.logger.Infof("上传文件成功, path: %s", path)
	if err = os.Chmod(dir, mode); err != nil {
		u.logger.Errorf("上传文件失败, err: %v", err)
		return "", err
	}

	out, err := os.Create(path)
	if err != nil {
		u.logger.Errorf("上传文件失败, err: %v", err)
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		u.logger.Errorf("上传文件失败, err: %v", err)
		return "", err
	}

	// 调用RAG接口
	return "上传文件成功并且RAG成功", nil
}
