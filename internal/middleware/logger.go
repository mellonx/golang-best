package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mellonx/golang-best/pkg/logger"
	"github.com/mellonx/golang-best/pkg/utils"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 生成唯一的Request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 设置Request ID到上下文
		ctx := logger.SetRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// 设置Response Header中的Request ID
		c.Header("X-Request-ID", requestID)

		// 获取请求信息
		referer := utils.GetReferer(c)
		clientIP := utils.GetClientIP(c)
		userAgent := utils.GetUserAgent(c)
		params := utils.ExtractRequestParams(c)

		// 打印请求开始日志
		logger.InfoCtx(ctx, "[Request Started]",
			"method", method,
			"path", path,
			"client_ip", clientIP,
			"referer", referer,
			"user_agent", userAgent,
			"params", params,
		)

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 打印请求完成日志
		logger.InfoCtx(ctx, "[Request Completed]",
			"status", statusCode,
			"method", method,
			"path", path,
			"latency", latency.String(),
			"client_ip", clientIP,
		)
	}
}
