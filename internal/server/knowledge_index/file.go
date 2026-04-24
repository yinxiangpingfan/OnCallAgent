package knowledgeindex

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
	"github.com/sirupsen/logrus"
)

type fileUploader struct {
	logger *logrus.Logger
	r      compose.Runnable[document.Source, bool]
}

type FileUploader interface {
	Upload(ctx context.Context, file *multipart.FileHeader, path string) (string, error)
}

func NewFileUploaderServer(log *logrus.Logger) FileUploader {
	return &fileUploader{logger: log}
}

func (u *fileUploader) Upload(ctx context.Context, file *multipart.FileHeader, path string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	path += file.Filename

	var mode os.FileMode = 0o750

	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, mode); err != nil {
		u.logger.Errorf("上传文件失败, err: %v", err)
		return "", err
	}
	if err = os.Chmod(dir, mode); err != nil {
		u.logger.Errorf("上传文件失败, err: %v", err)
		return "", err
	}

	out, err := os.Create(path)
	if err != nil {
		u.logger.Errorf("上传文件失败, err: %v", err)
		return "", err
	}

	_, err = io.Copy(out, src)
	if err != nil {
		u.logger.Errorf("写入文件失败, err: %v", err)
		return "", err
	}

	// 显式关闭文件，确保数据落盘后再调用 RAG
	if err = out.Close(); err != nil {
		u.logger.Errorf("关闭文件失败, err: %v", err)
		return "", err
	}

	// 调用 RAG 接口；若失败则回滚（删除已写入的文件）
	_, err = u.r.Invoke(ctx, document.Source{
		URI: "./docs/" + file.Filename,
	})
	if err != nil {
		u.logger.Errorf("RAG 失败，开始回滚删除文件 %s, err: %v", path, err)
		if removeErr := os.Remove(path); removeErr != nil {
			u.logger.Errorf("回滚删除文件失败, path: %s, err: %v", path, removeErr)
		}
		return "", fmt.Errorf("RAG 失败，文件已回滚: %w", err)
	}
	return "上传文件成功并且RAG成功", nil
}
