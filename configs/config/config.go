// Package config 提供应用程序配置管理功能
package config

// 保留原有导入和结构体定义
// 仅修改包声明为config
import (
	"fmt"
	"gin-center/pkg/security/useJwt"

	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// instance 保存配置单例
	instance *GlobalConfig
	// once 确保配置只被加载一次
	once sync.Once
	// validate 用于配置验证
	validate *validator.Validate
)

// BaseConfig 定义基础配置项
type BaseConfig struct {
	// Name 应用名称
	Name string `mapstructure:"name" validate:"required"`
	// Env 运行环境，可选值：development/production/testing
	Env string `mapstructure:"env" validate:"required,oneof=development production testing"`
	// Port 服务端口号
	Port int `mapstructure:"port" validate:"required"`
	// LogLevel 日志级别，可选值：debug/info/warn/error
	LogLevel string `mapstructure:"log_level" validate:"required,oneof=debug info warn error"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	// Driver 数据库驱动名称
	Driver string `mapstructure:"driver" validate:"required"`
	// Host 数据库主机地址
	Host string `mapstructure:"host" validate:"required"`
	// Port 数据库端口
	Port int `mapstructure:"port" validate:"required"`
	// Username 数据库用户名
	Username string `mapstructure:"username" validate:"required"`
	// Password 数据库密码
	Password string `mapstructure:"password"`
	// DBName 数据库名称
	DBName string `mapstructure:"dbname" validate:"required"`
	// Charset 字符集
	Charset string `mapstructure:"charset" validate:"required"`
	// ParseTime 是否解析时间
	ParseTime bool `mapstructure:"parse_time"`
	// Location 时区设置
	Location string `mapstructure:"location"`
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns" validate:"required"`
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `mapstructure:"max_open_conns" validate:"required"`
	// MaxLifetime 连接最大生命周期
	MaxLifetime time.Duration `mapstructure:"max_lifetime"`
	// MaxIdleTime 空闲连接最大生命周期
	MaxIdleTime time.Duration `mapstructure:"max_idle_time"`
}

// 基础Redis服务配置
type RedisConfig struct {
	// Host Redis主机地址
	Host string `mapstructure:"host"`
	// Port Redis端口
	Port int `mapstructure:"port"`
	// Password Redis密码
	Password string `mapstructure:"password"`
	// DB Redis数据库索引
	DB int `mapstructure:"db"`
	// PoolSize 连接池大小
	PoolSize int `mapstructure:"pool_size"`
	// MinIdleConns 最小空闲连接数
	MinIdleConns int `mapstructure:"min_idle_conns"`
	// MaxConnAge 连接最大年龄
	MaxConnAge time.Duration `mapstructure:"max_conn_age"`
	// IdleTimeout 空闲超时时间
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
	// Wait 连接池满时是否等待
	Wait bool `mapstructure:"wait"`
	// MaxRetries 最大重试次数
	MaxRetries int `mapstructure:"max_retries"`
	// MinRetryDelay 最小重试延迟
	MinRetryDelay time.Duration `mapstructure:"min_retry_delay"`
	// MaxRetryDelay 最大重试延迟
	MaxRetryDelay time.Duration `mapstructure:"max_retry_delay"`
	// MinRetryBackoff 最小重试退避时间
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`
	// MaxRetryBackoff 最大重试退避时间
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`
}

// LogConfig 日志配置
type LogConfig struct {
	// Level 日志级别
	Level string `mapstructure:"level" validate:"required"`
	// Filename 日志文件名
	Filename string `mapstructure:"filename" validate:"required"`
	// MaxSize 单个日志文件最大尺寸，单位MB
	MaxSize int `mapstructure:"max_size" validate:"required"`
	// MaxBackups 最大保留的日志文件数
	MaxBackups int `mapstructure:"max_backups" validate:"required"`
	// MaxAge 日志文件保留天数
	MaxAge int `mapstructure:"max_age" validate:"required"`
	// Compress 是否压缩历史日志
	Compress bool `mapstructure:"compress"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	// CORS 跨域配置
	CORS CORSConfig `mapstructure:"cors"`
	// Port 服务器监听端口
	Port int `mapstructure:"port" validate:"required"`
	// ReadTimeout 读取超时时间
	ReadTimeout time.Duration `mapstructure:"read_timeout" default:"5s"`
	// WriteTimeout 写入超时时间
	WriteTimeout time.Duration `mapstructure:"write_timeout" default:"10s"`
	// IdleTimeout 空闲连接超时时间
	IdleTimeout time.Duration `mapstructure:"idle_timeout" default:"15s"`
	// ReadHeaderTimeout 读取请求头超时时间
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout" default:"3s"`
	// ShutdownTimeout 优雅关闭超时时间
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" default:"30s"`
	// MaxHeaderBytes 请求头最大字节数
	MaxHeaderBytes int `mapstructure:"max_header_bytes" default:"1048576"`
	// MaxConnsPerIP 每个IP的最大连接数
	MaxConnsPerIP int `mapstructure:"max_conns_per_ip"`
	// MaxRequestsPerConn 每个连接的最大请求数
	MaxRequestsPerConn int `mapstructure:"max_requests_per_conn"`
}

// GlobalConfig 应用程序总配置结构
type GlobalConfig struct {
	Logger *zap.Logger
	mu     sync.RWMutex

	App      AppConfig        `mapstructure:"app"`
	Server   ServerConfig     `mapstructure:"server"`
	Database DatabaseConfig   `mapstructure:"database"`
	Redis    RedisConfig      `mapstructure:"redis"`
	Log      LogConfig        `mapstructure:"log"`
	JWT      useJwt.JWTConfig `mapstructure:"jwt"`
}

// 调整AppConfig结构体映射方式
type AppConfig struct {
	Name       string `mapstructure:"name" validate:"required"`
	Env        string `mapstructure:"env" validate:"required"`
	Port       int    `mapstructure:"port" validate:"required"`
	LogLevel   string `mapstructure:"log_level"`
	Version    string `mapstructure:"version"`
	Host       string `mapstructure:"host"`
	SuperAdmin string `mapstructure:"super_admin"` // 新增超级管理员配置字段
}

// LoadConfig 加载应用程序配置
// configPath: 配置文件所在目录路径
func LoadConfig(configPath string) (*GlobalConfig, error) {
	var err error
	once.Do(func() {
		instance = &GlobalConfig{}
		validate = validator.New()

		// 初始化viper配置
		if err = initViperConfig(configPath); err != nil {
			return
		}

		// 加载环境特定配置
		if err = loadEnvConfig(configPath); err != nil {
			log.Printf("加载环境配置失败: %v\n", err)
			return
		}

		// 解析并验证配置
		if err = parseConfig(); err != nil {
			return
		}

		// 设置配置文件监听
		setupConfigWatcher()
	})

	if err != nil {
		return nil, fmt.Errorf("配置加载失败: %w", err)
	}
	return instance, nil
}

// initViperConfig 初始化Viper配置管理器
// configPath: 配置文件所在目录路径
func initViperConfig(configPath string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}
	return nil
}

// loadEnvConfig 加载环境特定的配置文件
// configPath: 配置文件所在目录路径
func loadEnvConfig(configPath string) error {
	env := viper.GetString("APP_ENV")
	if env == "" {
		env = "dev"
	}

	envConfigPath := filepath.Join(configPath, fmt.Sprintf("%s.yaml", env))
	if _, err := os.Stat(envConfigPath); err == nil {
		viper.SetConfigName(env)
		if err := viper.MergeInConfig(); err != nil {
			return fmt.Errorf("合并环境配置失败: %w", err)
		}
	}
	return nil
}

// parseConfig 解析并验证配置
// 将配置文件内容解析到结构体中，并进行验证
func parseConfig() error {
	if err := viper.Unmarshal(instance); err != nil {
		return fmt.Errorf("配置解析失败: %w", err)
	}

	if err := validate.Struct(instance); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}
	return nil
}

// setupConfigWatcher 设置配置文件变更监听
// 监听配置文件变更，自动重新加载配置
func setupConfigWatcher() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		instance.mu.Lock()
		defer instance.mu.Unlock()

		if err := parseConfig(); err != nil {
			log.Printf("配置重新加载失败: %v\n", err)
			return
		}
		log.Printf("配置已更新，文件: %s\n", e.Name)
	})
}

// GetConfig 获取配置实例
// 返回值: 全局配置实例
func GetConfig() *GlobalConfig {
	return instance
}

// GetConfigByKey 根据键名获取配置项
// 返回值: (配置项, 是否存在)
func (c *GlobalConfig) GetConfigByKey(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	configs := map[string]any{
		"app":      &c.App,
		"server":   &c.Server,
		"database": &c.Database,
		"redis":    &c.Redis,
		"log":      &c.Log,
		"jwt":      &c.JWT,
	}

	config, exists := configs[key]
	return config, exists
}

// IsDevelopment 判断是否为开发环境
// 返回值: 如果当前环境为development则返回true
func (c *GlobalConfig) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction 判断是否为生产环境
// 返回值: 如果当前环境为production则返回true
func (c *GlobalConfig) IsProduction() bool {
	return c.App.Env == "production"
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}
