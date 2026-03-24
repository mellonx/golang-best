package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int
	Mode         string
	ReadTimeout  int
	WriteTimeout int
	CORS         CORSConfig
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string
	Format string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

// Load 加载配置
func Load() (*Config, error) {
	// 设置配置文件名和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 从环境变量读取
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// 配置文件不存在时使用默认值
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置默认配置
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 60)
	viper.SetDefault("server.write_timeout", 60)

	// CORS默认配置
	viper.SetDefault("server.cors.allow_origins", []string{"*"})
	viper.SetDefault("server.cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("server.cors.allow_headers", []string{"Origin", "Content-Type", "Authorization"})
	viper.SetDefault("server.cors.allow_credentials", true)
	viper.SetDefault("server.cors.max_age", 86400)

	// 数据库默认配置
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "app")
	viper.SetDefault("database.sslmode", "disable")

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// JWT默认配置
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expire_time", "24h")
}
