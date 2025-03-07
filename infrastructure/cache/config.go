package cache

import "time"

type Config struct {
	Redis RedisConfig
}

// 缓存专用Redis实例配置
type RedisConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Password        string        `yaml:"password"`
	DB              int           `yaml:"db"`
	PoolSize        int           `yaml:"pool_size"`
	MinIdleConns    int           `yaml:"min_idle_conns"`
	MaxConnAge      time.Duration `yaml:"max_conn_age"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	MaxRetries      int           `yaml:"max_retries"`
	MinRetryBackoff time.Duration `yaml:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `yaml:"max_retry_backoff"`
}
