package main

import (
	"log"

	"github.com/mellonx/golang-best/internal/config"
	"github.com/mellonx/golang-best/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 启动API服务器
	if err := RunAPIServer(cfg); err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}
