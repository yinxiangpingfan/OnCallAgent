package handler

import (
	"OnCallAgent/internal/server/chatServer"

	"github.com/gin-gonic/gin"
)

type chatHandler struct {
	chat chatServer.ChatServer
}

type ChatHandler interface {
	Chat() gin.HandlerFunc
	ChatSream() gin.HandlerFunc
}

func NewChatHandler(chat chatServer.ChatServer) ChatHandler {
	return &chatHandler{chat: chat}
}

type ChatRequest struct {
	Question string `json:"question" binding:"required"`
	Id       string `json:"id" binding:"required"` // 会话id
}

func (c *chatHandler) Chat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var chatRequest ChatRequest
		if err := ctx.ShouldBindJSON(&chatRequest); err != nil {
			ctx.JSON(400, gin.H{"message": "invalid request"})
			return
		}
		msg, err := c.chat.Chat(ctx.Request.Context(), chatRequest.Question, chatRequest.Id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "对话失败",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": msg,
		})
	}
}

func (c *chatHandler) ChatSream() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
