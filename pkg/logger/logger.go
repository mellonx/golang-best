package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.SugaredLogger
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

	logger = zapLogger.Sugar()
	return nil
}

// Debug 调试日志
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof 格式化信息日志
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf 格式化错误日志
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
	os.Exit(1)
}

// Fatalf 格式化致命错误日志
func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
	os.Exit(1)
}

// Sync 同步日志
func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}
