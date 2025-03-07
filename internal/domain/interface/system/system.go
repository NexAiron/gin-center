package system_config

import (
	"gin-center/internal/types/system"
)

type SystemServiceInterface interface {
	GetSystemConfig() map[string]any
	UpdateSystemConfig(config system.SystemConfig) error
	GetSystemInfo() (*system.SystemConfig, error)
}
