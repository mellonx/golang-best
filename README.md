# Golang Best Practices

一个基于 Gin 框架的 Golang 项目最佳实践模板，采用标准的项目结构和清晰的分层架构。

## 项目特点

- **标准项目结构**：遵循 Go 社区最佳实践
- **分层架构**：Controller → Service → Repository 清晰分层
- **依赖注入**：便于测试和维护
- **配置管理**：使用 Viper 支持多种配置格式
- **日志系统**：集成 Zap 高性能日志库
- **优雅关闭**：支持优雅关闭 HTTP 服务器
- **CORS 支持**：内置跨域中间件
- **统一响应**：标准化的 API 响应格式

## 项目结构

```
golang-best/
├── cmd/                    # 应用入口
│   └── api/               # API服务入口
│       ├── main.go        # 主程序入口
│       └── server.go      # 服务器设置
├── configs/               # 配置文件
│   └── config.yaml       # 默认配置
├── internal/             # 私有应用代码
│   ├── config/          # 配置加载
│   ├── controllers/     # 控制器层（处理HTTP请求）
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   ├── repositories/    # 数据访问层
│   ├── services/        # 业务逻辑层
│   └── utils/           # 内部工具函数
├── pkg/                  # 公共库（可被外部引用）
│   ├── api/             # API相关工具
│   ├── logger/          # 日志封装
│   └── response/        # 统一响应格式
├── api/                  # API文档
├── docs/                 # 项目文档
├── scripts/              # 脚本文件
├── go.mod                # Go模块定义
└── README.md             # 项目说明
```

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL（或其他数据库）

### 安装依赖

```bash
go mod download
```

### 配置

复制配置文件并根据需要修改：

```bash
cp configs/config.yaml configs/config.local.yaml
```

### 运行

```bash
go run cmd/api/main.go
```

服务将在 `http://localhost:8080` 启动。

### 健康检查

```bash
curl http://localhost:8080/health
```

## API 示例

### 用户管理

```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'

# 获取用户列表
curl http://localhost:8080/api/v1/users

# 获取单个用户
curl http://localhost:8080/api/v1/users/1

# 更新用户
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Smith"
  }'

# 删除用户
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 架构说明

### 分层架构

1. **Controller 层**：负责处理 HTTP 请求和响应，参数验证
2. **Service 层**：业务逻辑处理，事务管理
3. **Repository 层**：数据访问，数据库操作

### 依赖流向

```
Controller → Service → Repository
```

每一层只依赖下一层的接口，便于单元测试和模块替换。

## 配置说明

配置文件采用 YAML 格式，支持环境变量覆盖：

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `server.port` | 服务端口 | 8080 |
| `server.mode` | 运行模式 | debug |
| `database.host` | 数据库主机 | localhost |
| `database.port` | 数据库端口 | 5432 |
| `log.level` | 日志级别 | info |

## 开发指南

### 添加新的 API

1. 在 `internal/models/` 定义数据模型
2. 在 `internal/repositories/` 创建 Repository 接口和实现
3. 在 `internal/services/` 创建 Service 接口和实现
4. 在 `internal/controllers/` 创建 Controller
5. 在 `cmd/api/server.go` 注册路由

### 运行测试

```bash
go test ./... -v
```

### 构建二进制文件

```bash
go build -o bin/api cmd/api/main.go
```

## 最佳实践

1. **错误处理**：统一错误响应格式，记录详细错误日志
2. **日志规范**：使用结构化日志，便于查询和分析
3. **配置管理**：敏感信息使用环境变量，不提交到代码库
4. **代码组织**：按功能模块组织，保持单一职责原则
5. **依赖注入**：避免全局变量，便于测试

## 技术栈

- **Web 框架**：[Gin](https://github.com/gin-gonic/gin)
- **配置管理**：[Viper](https://github.com/spf13/viper)
- **日志库**：[Zap](https://github.com/uber-go/zap)
- **ORM**：[GORM](https://gorm.io/)

## License

MIT License
