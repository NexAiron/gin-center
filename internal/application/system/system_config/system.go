package systemConfig

import (
	"context"
	"fmt"

	"runtime"
	"time"

	"gin-center/infrastructure/zaplogger"
	use_Baseservice "gin-center/internal/application"
	"gin-center/internal/types/system"

	"github.com/go-redis/redis/v8"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type SystemService struct {
	*use_Baseservice.BaseService
	config *viper.Viper
	redis  *redis.Client
}

func NewSystemService(redis *redis.Client, logger *zaplogger.ServiceLogger) *SystemService {
	baseService := use_Baseservice.NewBaseService(&use_Baseservice.BaseServiceConfig{
		Logger: logger,
	})
	return &SystemService{
		BaseService: baseService,
		config:      viper.GetViper(),
		redis:       redis,
	}
}
func (s *SystemService) getConfigValue(key string, defaultValue interface{}) interface{} {
	if !s.config.IsSet(key) {
		s.BaseService.Logger.LogWarn("配置项不存在，使用默认值", zap.Skip())
		return defaultValue
	}
	return s.config.Get(key)
}
func (s *SystemService) GetSystemConfig() map[string]interface{} {
	return map[string]interface{}{
		"system_name": s.getConfigValue("system.name", "gin-center"),
		"version":     s.getConfigValue("system.version", "1.0.0"),
		"debug_mode":  s.getConfigValue("system.debug", false),
		"server": map[string]interface{}{
			"host": s.getConfigValue("server.host", "localhost"),
			"port": s.getConfigValue("server.port", 8080),
		},
	}
}
func (s *SystemService) UpdateSystemConfig(config system.SystemConfig) error {
	s.BaseService.Logger.LogInfo("开始更新系统配置", zap.String("operation", "config_update"))
	s.config.Set("system.name", config.SystemName)
	s.config.Set("system.version", config.Version)
	s.config.Set("system.environment", config.Environment)
	s.config.Set("system.debug_mode", config.DebugMode)
	s.config.Set("system.log_level", config.LogLevel)
	s.config.Set("system.port", config.Port)
	s.config.Set("security.jwt_secret", config.JWTSecret)
	s.config.Set("security.jwt_expire", config.JWTExpire)
	s.config.Set("database.host", config.DBHost)
	s.config.Set("database.port", config.DBPort)
	s.config.Set("database.name", config.DBName)
	s.config.Set("database.user", config.DBUser)
	s.config.Set("database.password", config.DBPassword)
	s.config.Set("redis.host", config.RedisHost)
	s.config.Set("redis.port", config.RedisPort)
	s.config.Set("redis.db", config.RedisDB)
	s.config.Set("redis.password", config.RedisPassword)
	if err := s.validateConfig(config); err != nil {
		s.BaseService.Logger.LogError("配置验证失败", zap.Skip(), zap.Error(err))
		return fmt.Errorf("配置验证失败: %w", err)
	}
	if err := s.config.WriteConfig(); err != nil {
		s.BaseService.Logger.LogError("写入配置文件失败", zap.Skip(), zap.Error(err))
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	s.BaseService.Logger.LogInfo("系统配置更新成功", zap.String("operation", "config_update"))
	return nil
}
func (s *SystemService) validateConfig(config system.SystemConfig) error {
	if len(config.SystemName) == 0 {
		return fmt.Errorf("系统名称不能为空")
	}
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("端口号必须在1-65535之间")
	}
	if len(config.JWTSecret) < 16 {
		return fmt.Errorf("JWT密钥长度必须大于等于16位")
	}
	return nil
}
func (s *SystemService) GetSystemHealth() map[string]interface{} {
	metrics, err := s.GetSystemMetrics()
	status := "up"
	details := map[string]string{
		"database":   "up",
		"redis":      s.checkRedisHealth(),
		"jwt_auth":   "up",
		"log_system": "up",
	}

	if err != nil {
		status = "degraded"
		details["metrics"] = "unavailable"
		s.BaseService.Logger.LogError("获取系统指标失败", zap.Skip(), zap.Error(err))
	} else {
		if mem, ok := metrics["memory"].(map[string]interface{}); ok {
			if usedPercent, ok := mem["used_percent"].(float64); ok && usedPercent > 90 {
				status = "warning"
				details["memory"] = "high_usage"
				s.BaseService.Logger.LogWarn("内存使用率过高", zap.Skip(), zap.Float64("used_percent", usedPercent))
			}
		}
		if cpu, ok := metrics["cpu"].(map[string]interface{}); ok {
			if usagePercent, ok := cpu["usage_percent"].(float64); ok && usagePercent > 80 {
				status = "warning"
				details["cpu"] = "high_load"
				s.BaseService.Logger.LogWarn("CPU负载过高", zap.Skip(), zap.Float64("usage_percent", usagePercent))
			}
		}
	}

	result := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
		"details":   details,
	}

	info, err := s.GetSystemInfo()
	if err == nil && info != nil {
		result["system_info"] = info
	}

	return result
}

func (s *SystemService) GetSystemInfo() (*system.SystemConfig, error) {
	metrics, err := s.GetSystemMetrics()
	if err != nil {
		s.BaseService.Logger.LogError("获取系统指标失败", zap.Skip(), zap.Error(err))
		return nil, err
	}

	info := &system.SystemConfig{}
	info.SystemName = s.getConfigValue("system.name", "gin-center").(string)
	info.Version = s.getConfigValue("system.version", "1.0.0").(string)
	info.App.Uptime = time.Now().Format(time.RFC3339)

	if mem, ok := metrics["memory"].(map[string]interface{}); ok {
		info.Stats.MemoryTotal = mem["total"].(uint64)
		info.Stats.MemoryUsed = mem["used"].(uint64)
		info.Stats.MemoryUsedPercent = mem["used_percent"].(float64)
	}

	if cpu, ok := metrics["cpu"].(map[string]interface{}); ok {
		info.Stats.CPUUsagePercent = cpu["usage_percent"].(float64)
	}

	return info, nil
}

func (s *SystemService) GetSystemMetrics() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		s.BaseService.Logger.LogError("获取内存信息失败", zap.Skip(), zap.Error(err))
		return nil, fmt.Errorf("获取内存信息失败: %w", err)
	}

	var cpuPercent []float64
	for i := 0; i < 3; i++ {
		cpuPercent, err = cpu.PercentWithContext(ctx, time.Second, false)
		if err == nil {
			break
		}
		s.BaseService.Logger.LogWarn("获取CPU信息失败,正在重试", zap.Skip(), zap.Error(err))
		time.Sleep(time.Second)
	}
	if err != nil {
		s.BaseService.Logger.LogError("获取CPU信息失败", zap.Skip(), zap.Error(err))
		return nil, fmt.Errorf("获取CPU信息失败: %w", err)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"memory": map[string]interface{}{
			"total":        vmStat.Total,
			"available":    vmStat.Available,
			"used":         vmStat.Used,
			"used_percent": vmStat.UsedPercent,
		},
		"cpu": map[string]interface{}{
			"usage_percent": cpuPercent[0],
		},
		"go_runtime": map[string]interface{}{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"num_gc":      m.NumGC,
			"goroutines":  runtime.NumGoroutine(),
		},
	}, nil
}
func (s *SystemService) checkRedisHealth() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	result := s.redis.Ping(ctx)
	if result.Err() != nil {
		s.BaseService.Logger.LogWarn("Redis连接异常", zap.Error(result.Err()))
		return "down"
	}
	return "up"
}
