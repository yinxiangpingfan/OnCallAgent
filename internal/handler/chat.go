package handler

import (
	"OnCallAgent/internal/server/chatServer"
	"fmt"

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
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		var chatRequest ChatRequest
		if err := ctx.ShouldBindJSON(&chatRequest); err != nil {
			ctx.JSON(400, gin.H{"message": "invalid request"})
			return
		}

		// 带缓冲 channel，避免 goroutine 在客户端断开后仍阻塞在写入
		ch := make(chan string, 8)
		done := make(chan struct{})

		go c.chat.ChatSream(ctx.Request.Context(), chatRequest.Question, chatRequest.Id, &ch, &done)

		for {
			// 每次循环前，先看一眼 done 有没有动静
			select {
			case <-done:
				return
			default:
			}
			t, ok := <-ch
			if !ok {
				// channel 已关闭，goroutine 正常结束
				ctx.SSEvent("message", "data: [DONE]\n\n")
				ctx.Writer.Flush()
				return
			}
			// SSE 数据格式要求（重要！）
			event := fmt.Sprintf("data: %v\n\n", t)
			// 发送数据到客户端
			ctx.SSEvent("message", event)
			ctx.Writer.Flush() // 立即刷新缓冲区
		}
	}
}
