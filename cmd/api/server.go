package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mellonx/golang-best/internal/config"
	"github.com/mellonx/golang-best/internal/controllers"
	"github.com/mellonx/golang-best/internal/middleware"
	"github.com/mellonx/golang-best/internal/repositories"
	"github.com/mellonx/golang-best/internal/services"
	"github.com/mellonx/golang-best/pkg/logger"
	"github.com/mellonx/golang-best/pkg/response"
	"gorm.io/gorm"
)

// RunAPIServer 启动API服务器
func RunAPIServer(cfg *config.Config) error {
	// 设置Gin运行模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库连接
	db, err := initDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}

	// 自动迁移数据库表
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 初始化路由
	router := setupRouter(cfg, db)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器（在goroutine中）
	go func() {
		logger.Info(fmt.Sprintf("Server starting on port %d", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("Server exited")
	return nil
}

// initDatabase 初始化数据库连接
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	// 这里使用GORM连接数据库
	// 实际项目中根据配置选择数据库驱动
	// 示例：连接PostgreSQL
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
	// 	cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	// return gorm.Open(postgres.Open(dsn), &gorm.Config{})

	logger.Info("Database initialization (mock for demo)")
	return nil, nil
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	// 在这里添加模型迁移
	// return db.AutoMigrate(&models.User{}, &models.Post{})
	return nil
}

// setupRouter 设置路由
func setupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.Server.CORS))

	// 初始化依赖
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	// API路由组
	api := router.Group("/api/v1")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.GET("", userController.List)
			users.GET("/:id", userController.GetByID)
			users.POST("", userController.Create)
			users.PUT("/:id", userController.Update)
			users.DELETE("/:id", userController.Delete)
		}
	}

	return router
}
