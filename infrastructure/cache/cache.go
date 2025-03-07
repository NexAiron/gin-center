// Package cache 提供了应用程序的缓存管理功能
// 实现了基于Redis的缓存操作，支持基本的缓存读写和批量操作
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	// ErrKeyNotFound 表示缓存中未找到指定的键
	ErrKeyNotFound = errors.New("key not found in cache")
)

// Cache 定义了缓存操作的接口
// 提供了基本的缓存读写功能
type Cache interface {
	// Get 获取缓存中的值
	Get(ctx context.Context, key string) (interface{}, error)
	// Set 设置缓存值，可指定过期时间
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Delete 删除缓存中的键
	Delete(ctx context.Context, key string) error
	// Unmarshal 将缓存数据反序列化到指定的结构体中
	Unmarshal(data interface{}, value interface{}) error
}

// RedisCache 实现了Cache接口的Redis缓存结构
type RedisCache struct {
	// client Redis客户端实例
	client *redis.Client
}

// NewRedisCache 创建一个新的Redis缓存实例
// client: Redis客户端实例
// 返回实现了Cache接口的RedisCache实例
func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Unmarshal 实现 Cache 接口的 Unmarshal 方法
// 使用 json 包将缓存数据反序列化到指定的结构体中
func (c *RedisCache) Unmarshal(data interface{}, value interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}
