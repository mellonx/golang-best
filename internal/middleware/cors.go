package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mellonx/golang-best/internal/config"
	"github.com/mellonx/golang-best/pkg/logger"
)

// CORS 跨域中间件
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		// 检查是否允许该源
		allowed := false
		for _, o := range cfg.AllowOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if !allowed {
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", joinStrings(cfg.AllowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", joinStrings(cfg.AllowHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", joinStrings(cfg.ExposeHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", boolToString(cfg.AllowCredentials))
		c.Header("Access-Control-Max-Age", intToString(cfg.MaxAge))

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", "error", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

// 辅助函数
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func intToString(i int) string {
	return string(rune(i + '0'))
}
