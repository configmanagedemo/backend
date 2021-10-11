package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.DebugLevel)
	logger.Out = os.Stdout
}

// LogMiddleware 日志处理中间件
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		// 耗时
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqURI := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqURI,
		)
	}
}

// INFO info日志
func INFO(message string) {
	logger.Infof("| INFO | %s |", message)
}

// DEBUG debug日志
func DEBUG(message string) {
	logger.Debugf("| DEBUG | %s |", message)
}

// WARN debug日志
func WARN(message string) {
	logger.Debugf("| WARN | %s |", message)
}

// ERROR debug日志
func ERROR(message string) {
	logger.Debugf("| ERROR | %s |", message)
}
