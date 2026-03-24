package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 上下文键类型
type contextKey string

const (
	// RequestIDKey 请求ID上下文键
	RequestIDKey contextKey = "requestId"
)

var (
	baseLogger *zap.SugaredLogger
)

// Init 初始化日志
func Init(level, format string) error {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 设置日志级别
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	zapLogger, err := config.Build()
	if err != nil {
		return err
	}

	baseLogger = zapLogger.Sugar()
	return nil
}

// WithRequestID 创建带有requestId的日志上下文
func WithRequestID(requestID string) *zap.SugaredLogger {
	return baseLogger.With("request_id", requestID)
}

// WithContext 从上下文获取带requestId的logger
func WithContext(ctx context.Context) *zap.SugaredLogger {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		return WithRequestID(requestID)
	}
	return baseLogger
}

// GetRequestID 从上下文获取requestId
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// SetRequestID 设置requestId到上下文
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// Debug 调试日志
func Debug(args ...interface{}) {
	baseLogger.Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(template string, args ...interface{}) {
	baseLogger.Debugf(template, args...)
}

// DebugCtx 带上下文的调试日志
func DebugCtx(ctx context.Context, args ...interface{}) {
	WithContext(ctx).Debug(args...)
}

// DebugfCtx 带上下文的格式化调试日志
func DebugfCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Debugf(template, args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	baseLogger.Info(args...)
}

// Infof 格式化信息日志
func Infof(template string, args ...interface{}) {
	baseLogger.Infof(template, args...)
}

// InfoCtx 带上下文的信息日志
func InfoCtx(ctx context.Context, args ...interface{}) {
	WithContext(ctx).Info(args...)
}

// InfofCtx 带上下文的格式化信息日志
func InfofCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Infof(template, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	baseLogger.Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(template string, args ...interface{}) {
	baseLogger.Warnf(template, args...)
}

// WarnCtx 带上下文的警告日志
func WarnCtx(ctx context.Context, args ...interface{}) {
	WithContext(ctx).Warn(args...)
}

// WarnfCtx 带上下文的格式化警告日志
func WarnfCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Warnf(template, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	baseLogger.Error(args...)
}

// Errorf 格式化错误日志
func Errorf(template string, args ...interface{}) {
	baseLogger.Errorf(template, args...)
}

// ErrorCtx 带上下文的错误日志
func ErrorCtx(ctx context.Context, args ...interface{}) {
	WithContext(ctx).Error(args...)
}

// ErrorfCtx 带上下文的格式化错误日志
func ErrorfCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Errorf(template, args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	baseLogger.Fatal(args...)
	os.Exit(1)
}

// Fatalf 格式化致命错误日志
func Fatalf(template string, args ...interface{}) {
	baseLogger.Fatalf(template, args...)
	os.Exit(1)
}

// Sync 同步日志
func Sync() error {
	if baseLogger != nil {
		return baseLogger.Sync()
	}
	return nil
}
