package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitLogger 初始化日志记录器
func InitLogger(level string, path string) *logrus.Logger {
	logger := logrus.New() //创建一个logger实例
	switch level {
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	default:
		logger.SetLevel(logrus.ErrorLevel)
	}
	//普通日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,                      //强制不显示颜色
		TimestampFormat: "2006-01-02 15:04:05.000", // 显示ms
	})
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(file) //设置日志文件
	return logger
}
