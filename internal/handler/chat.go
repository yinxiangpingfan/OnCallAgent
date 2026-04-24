package handler

import "github.com/gin-gonic/gin"

type chatHandler struct {
}

type ChatHandler interface {
	Chat() gin.HandlerFunc
	ChatSream() gin.HandlerFunc
}

func NewChatHandler() ChatHandler {
	return &chatHandler{}
}

type ChatRequest struct {
	Question string `json:"question"`
}

func (c *chatHandler) Chat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var chatRequest ChatRequest
		if err := ctx.ShouldBindJSON(&chatRequest); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

	}
}

func (c *chatHandler) ChatSream() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
