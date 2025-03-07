package system

type SystemConfig struct {
	SystemName  string `json:"system_name" binding:"required"`
	Version     string `json:"version" binding:"required"`
	Environment string `json:"environment" binding:"required"`
	DebugMode   bool   `json:"debug_mode"`
	LogLevel    string `json:"log_level" binding:"required"`
	Port        int    `json:"port" binding:"required"`
	JWTSecret   string `json:"jwt_secret" binding:"required"`
	JWTExpire   int    `json:"jwt_expire" binding:"required"`
	Stats       struct {
		MemoryTotal       uint64  `json:"memory_total"`
		MemoryUsed        uint64  `json:"memory_used"`
		MemoryUsedPercent float64 `json:"memory_used_percent"`
		CPUUsagePercent   float64 `json:"cpu_usage_percent"`
	} `json:"stats"`
	DBHost        string `json:"db_host" binding:"required"`
	DBPort        int    `json:"db_port" binding:"required"`
	DBName        string `json:"db_name" binding:"required"`
	DBUser        string `json:"db_user" binding:"required"`
	DBPassword    string `json:"db_password" binding:"required"`
	RedisHost     string `json:"redis_host" binding:"required"`
	RedisPort     int    `json:"redis_port" binding:"required"`
	RedisDB       int    `json:"redis_db"`
	RedisPassword string `json:"redis_password"`
	App struct {
		Uptime string `json:"uptime"`
	} `json:"app"`
}
