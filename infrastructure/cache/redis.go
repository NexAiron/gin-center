// Package cache 提供了Redis缓存服务的初始化和管理功能
package cache

import (
	"context"
	"fmt"
	"gin-center/configs/config"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient 全局Redis客户端实例
// 用于提供对Redis服务器的连接和操作能力
var RedisClient *redis.Client

// InitRedis 初始化Redis客户端连接
func InitRedis(config *config.GlobalConfig) error {
	// 使用配置参数创建Redis客户端实例
	RedisClient = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password:        config.Redis.Password,
		DB:              config.Redis.DB,
		PoolSize:        config.Redis.PoolSize,
		MinIdleConns:    config.Redis.MinIdleConns,
		MaxConnAge:      config.Redis.MaxConnAge,
		IdleTimeout:     config.Redis.IdleTimeout,
		MaxRetries:      config.Redis.MaxRetries,
		MinRetryBackoff: config.Redis.MinRetryBackoff,
		MaxRetryBackoff: config.Redis.MaxRetryBackoff,
	})

	// 创建带超时的上下文，确保连接验证不会无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 验证Redis连接是否可用
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis连接失败: %v", err)
	}

	return nil
}

// CloseRedis 关闭Redis客户端连接
// 在应用程序退出时调用，确保资源被正确释放
// 返回error: 如果关闭连接时发生错误则返回错误信息
func CloseRedis() error {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			return fmt.Errorf("redis关闭失败: %v", err)
		}
	}
	return nil
}
