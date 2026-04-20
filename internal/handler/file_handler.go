package handler

import (
	errorsx "OnCallAgent/pkg/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileUploader interface {
	Upload(filepath string, data []byte) (string, error)
}

type Handler struct {
	uploader FileUploader
}

func NewHandler(uploader FileUploader) *Handler {
	return &Handler{uploader: uploader}
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "invalid request: no file provided",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "invalid request: cannot open file",
		})
		return
	}
	defer src.Close()

	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "invalid request: cannot read file",
		})
		return
	}

	filename := file.Filename
	url, err := h.uploader.Upload(filename, data)
	if err != nil {
		if errors.Is(err, errorsx.ErrFileTooLarge) {
			c.JSON(http.StatusOK, gin.H{
				"code": 413,
				"msg":  "file too large",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "internal error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": url,
	})
}
