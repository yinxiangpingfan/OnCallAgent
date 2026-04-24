package router

import (
	"OnCallAgent/internal/handler"
	knowledgeindex "OnCallAgent/internal/server/knowledge_index"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitRouter(r *gin.Engine, loger *logrus.Logger) {
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
	uploder := knowledgeindex.NewFileUploaderServer(loger)
	uploderHandler := handler.NewFileUploader("./docs/", uploder)
	r.POST("/upload", uploderHandler.Upload())
	//对话

}
