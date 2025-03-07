package zaplogger

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ServiceLogger struct {
	logger *zap.Logger
}

// LoggerConfig 日志配置项
type LoggerConfig struct {
	Level    string
	Filename string
}

// NewLogger 创建一个新的日志记录器
func NewServiceLogger() *ServiceLogger {
	// 创建默认配置
	cfg := LoggerConfig{
		Level: "info", // 默认日志级别为 info
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 设置时间格式为 ISO8601
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 级别显示为大写

	// 解析日志级别
	logLevel := zap.NewAtomicLevel()
	if err := logLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		logLevel.SetLevel(zap.InfoLevel) // 解析失败时回退到 Info 级别
	}

	// 创建核心组件
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON 编码器
		zapcore.AddSync(os.Stdout),            // 输出到标准输出
		logLevel,                              // 日志级别
	)

	// 创建 zap.Logger 并添加调用方信息
	zapLogger := zap.New(core, zap.AddCaller())

	return &ServiceLogger{logger: zapLogger}
}

// LogInfo 记录一条信息日志
func (l *ServiceLogger) LogInfo(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// LogError 记录一条错误日志
func (l *ServiceLogger) LogError(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// LogWarn 记录一条警告日志
func (l *ServiceLogger) LogWarn(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// LogDebug 记录一条调试日志
func (l *ServiceLogger) LogDebug(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// LogFatal 记录一条致命错误日志
func (l *ServiceLogger) LogFatal(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// Sync 同步日志缓冲区
func (l *ServiceLogger) Sync() error {
	return l.logger.Sync()
}

// With 添加字段到日志记录器
func (l *ServiceLogger) With(fields ...zap.Field) *ServiceLogger {
	return &ServiceLogger{
		logger: l.logger.With(fields...),
	}
}

// 包级错误定义
var ErrInitLoggerFailed = errors.New("日志初始化失败")
