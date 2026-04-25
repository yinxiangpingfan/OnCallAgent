package router

import (
	"OnCallAgent/internal/handler"
	"OnCallAgent/internal/server/ai/agent/chat"
	"OnCallAgent/internal/server/chatServer"
	knowledgeindex "OnCallAgent/internal/server/knowledge_index"
	"OnCallAgent/internal/server/plan"
	"OnCallAgent/pkg/config"
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitRouter(ctx context.Context, r *gin.Engine, loger *logrus.Logger, config *config.Config, runner compose.Runnable[document.Source, bool], runnerChat compose.Runnable[*chat.UserMessage, *schema.Message], model *openai.ChatModel, retriever *qdrant_retriever.Retriever) {
	//cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//文件上传
	uploder := knowledgeindex.NewFileUploaderServer(loger, runner)
	uploderHandler := handler.NewFileUploader("./docs/", uploder)
	r.POST("/upload", uploderHandler.Upload())
	//对话
	chater := chatServer.NewChatServer(loger, runnerChat)
	chaterHandler := handler.NewChatHandler(chater)
	r.POST("/chat", chaterHandler.Chat())
	r.POST("/chatStream", chaterHandler.ChatSream())
	//运维
	planer := plan.NewPlanServer(*config, model, loger, retriever)
	planerH := handler.NewPlanHandler(planer)
	r.GET("/plan", planerH.Plan())
}
