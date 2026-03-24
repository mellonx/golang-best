package utils

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

// 需要脱敏的字段名（不区分大小写）
var sensitiveFields = map[string]bool{
	"password":        true,
	"passwd":          true,
	"pwd":             true,
	"secret":          true,
	"token":           true,
	"access_token":    true,
	"refresh_token":   true,
	"api_key":         true,
	"apikey":          true,
	"authorization":   true,
	"credit_card":     true,
	"creditcard":      true,
	"ssn":             true,
	"social_security": true,
	"private_key":     true,
	"privatekey":      true,
}

// SensitiveMask 敏感信息的掩码
const SensitiveMask = "******"

// SanitizeMap 对map中的敏感字段进行脱敏
func SanitizeMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		lowerKey := strings.ToLower(key)
		if sensitiveFields[lowerKey] {
			result[key] = SensitiveMask
		} else {
			switch v := value.(type) {
			case map[string]interface{}:
				result[key] = SanitizeMap(v)
			case []interface{}:
				result[key] = sanitizeSlice(v)
			default:
				result[key] = value
			}
		}
	}
	return result
}

// sanitizeSlice 对slice中的元素进行脱敏
func sanitizeSlice(data []interface{}) []interface{} {
	result := make([]interface{}, len(data))
	for i, item := range data {
		switch v := item.(type) {
		case map[string]interface{}:
			result[i] = SanitizeMap(v)
		case []interface{}:
			result[i] = sanitizeSlice(v)
		default:
			result[i] = v
		}
	}
	return result
}

// SanitizeString 对JSON字符串中的敏感字段进行脱敏
func SanitizeString(jsonStr string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// 如果解析失败，返回原字符串
		return jsonStr
	}
	sanitized := SanitizeMap(data)
	result, err := json.Marshal(sanitized)
	if err != nil {
		return jsonStr
	}
	return string(result)
}

// ExtractRequestParams 提取请求参数（包含脱敏处理）
func ExtractRequestParams(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})

	// 提取Query参数
	queryParams := c.Request.URL.Query()
	if len(queryParams) > 0 {
		queryMap := make(map[string]interface{})
		for key, values := range queryParams {
			if len(values) == 1 {
				queryMap[key] = values[0]
			} else {
				queryMap[key] = values
			}
		}
		params["query"] = SanitizeMap(queryMap)
	}

	// 提取Post表单参数
	if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
		contentType := c.GetHeader("Content-Type")

		// 处理application/x-www-form-urlencoded
		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			_ = c.Request.ParseForm()
			formMap := make(map[string]interface{})
			for key, values := range c.Request.PostForm {
				if len(values) == 1 {
					formMap[key] = values[0]
				} else {
					formMap[key] = values
				}
			}
			if len(formMap) > 0 {
				params["form"] = SanitizeMap(formMap)
			}
		}

		// 处理application/json
		if strings.Contains(contentType, "application/json") {
			var jsonBody map[string]interface{}
			if err := c.ShouldBind(&jsonBody); err == nil && len(jsonBody) > 0 {
				params["body"] = SanitizeMap(jsonBody)
			}
		}
	}

	// 提取路径参数
	pathParams := c.Params
	if len(pathParams) > 0 {
		pathMap := make(map[string]interface{})
		for _, param := range pathParams {
			pathMap[param.Key] = param.Value
		}
		params["path"] = pathMap
	}

	return params
}

// GetReferer 获取请求的Referer
func GetReferer(c *gin.Context) string {
	return c.GetHeader("Referer")
}

// GetUserAgent 获取User-Agent
func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

// GetClientIP 获取客户端IP
func GetClientIP(c *gin.Context) string {
	// 尝试从X-Forwarded-For获取
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 尝试从X-Real-IP获取
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 使用ClientIP方法
	return c.ClientIP()
}
