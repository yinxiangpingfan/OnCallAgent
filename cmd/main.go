package main

import (
	"OnCallAgent/internal/router"
	"OnCallAgent/pkg/log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志记录器
	log := log.InitLogger("info", "log/OnCallAgent.log")
	r := gin.Default()
	// 初始化
	router.InitRouter(r)
}
