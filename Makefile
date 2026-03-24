.PHONY: help build run test clean

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 构建应用
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go

run: ## 运行应用
	@echo "Running..."
	@go run cmd/api/main.go

test: ## 运行测试
	@echo "Testing..."
	@go test ./... -v -cover

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

clean: ## 清理构建文件
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

deps: ## 安装依赖
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

lint: ## 运行代码检查
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## 格式化代码
	@echo "Formatting code..."
	@go fmt ./...

docker-build: ## 构建Docker镜像
	@echo "Building Docker image..."
	@docker build -t golang-best:latest .

docker-run: ## 运行Docker容器
	@echo "Running Docker container..."
	@docker run -p 8080:8080 golang-best:latest
