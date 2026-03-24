package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mellonx/golang-best/pkg/logger"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logger.Info("HTTP Request",
			"status", statusCode,
			"method", method,
			"path", path,
			"latency", latency.String(),
			"ip", clientIP,
		)
	}
}
