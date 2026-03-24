# pkg/utils 包说明

本包提供通用的工具函数，可被外部项目引用。

## 模块说明

### sanitizer.go - 数据脱敏工具

提供请求参数脱敏功能，用于日志输出时保护敏感信息。

#### 功能特性

- 自动识别敏感字段（password, token, api_key等）
- 支持JSON字符串脱敏
- 支持Map结构脱敏
- 支持嵌套数据结构脱敏

#### 使用示例

```go
import "github.com/mellonx/golang-best/pkg/utils"

// 脱敏Map
data := map[string]interface{}{
    "username": "john",
    "password": "secret123",
}
sanitized := utils.SanitizeMap(data)
// 输出: {"username": "john", "password": "******"}

// 提取请求参数（自动脱敏）
params := utils.ExtractRequestParams(c)
// 返回已脱敏的请求参数
```

#### 内置敏感字段列表

- password, passwd, pwd
- secret
- token, access_token, refresh_token
- api_key, apikey
- authorization
- credit_card, creditcard
- ssn, social_security
- private_key, privatekey

可以扩展 `sensitiveFields` map 来添加更多敏感字段。
