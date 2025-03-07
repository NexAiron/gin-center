package database

import (
	"context"
	"fmt"
	"time"

	"gin-center/infrastructure/zaplogger"

	"gin-center/configs/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxRetries     = 3
	retryDelay     = time.Second * 2
	connectTimeout = time.Second * 10
)

func InitDatabase(cfg *config.GlobalConfig) (*gorm.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("配置对象不能为空")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	gormLogger := NewZapGormLogger(zaplogger.NewServiceLogger())

	config := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	var db *gorm.DB
	var err error

	// 添加重试机制
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
		defer cancel()

		db, err = gorm.Open(mysql.Open(dsn), config)
		if err == nil {
			break
		}

		gormLogger.Info(ctx, "数据库连接失败，准备重试", zap.Error(err), zap.Int("重试次数", i+1))
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("数据库连接超时: %v", err)
		case <-time.After(retryDelay):
			continue
		}
	}

	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %v", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 验证数据库连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	ctx := context.Background()
	gormLogger.Info(ctx, "数据库连接成功",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.DBName),
	)
	return db, nil
}

// NewZapGormLogger 创建基于zap的gorm日志记录器
func NewZapGormLogger(serviceLogger *zaplogger.ServiceLogger) logger.Interface {
	return logger.New(
		&zapGormWriter{serviceLogger: serviceLogger},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

type zapGormWriter struct {
	serviceLogger *zaplogger.ServiceLogger
}

func (w *zapGormWriter) Printf(format string, args ...interface{}) {
	w.serviceLogger.LogInfo(fmt.Sprintf(format, args...))
}
