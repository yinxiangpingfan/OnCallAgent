package main

import (
	"OnCallAgent/internal/handler"
	"OnCallAgent/internal/server/knowledge_index"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	uploader := knowledge_index.NewUploader(knowledge_index.UploaderConfig{
		StoragePath: "./uploads",
		MaxSize:     10 << 20, // 10MB
	})

	h := handler.NewHandler(uploader)

	r.POST("/upload", h.UploadFile)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
