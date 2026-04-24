package handler

import (
	knowledgeindex "OnCallAgent/internal/server/knowledge_index"

	"github.com/gin-gonic/gin"
)

// FileUploader 上传文件接口

type FileUploader interface {
	Upload() gin.HandlerFunc
}

type fileUploader struct {
	uploadPath   string
	upLoadServer knowledgeindex.FileUploader
}

func NewFileUploader(uploadPath string, upLoadServer knowledgeindex.FileUploader) FileUploader {
	return &fileUploader{uploadPath: uploadPath, upLoadServer: upLoadServer}
}

func (u *fileUploader) Upload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "invalid request: no file provided",
			})
			return
		}
		msg, err := u.upLoadServer.Upload(ctx.Request.Context(), file, u.uploadPath)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "上传文件失败",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": msg,
		})
	}
}
