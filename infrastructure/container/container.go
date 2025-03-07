// Package container 提供了应用程序的依赖注入容器实现
package container

import (
	"context"
	"fmt"
	"gin-center/configs/config"
	"gin-center/infrastructure/cache"
	"gin-center/infrastructure/database"
	"gin-center/infrastructure/repository/admin"
	user_repo "gin-center/infrastructure/repository/user"
	"gin-center/infrastructure/zaplogger"
	AdminService "gin-center/internal/application/admin/service"
	systemService "gin-center/internal/application/system/system_service"
	user_service "gin-center/internal/application/user/service"
	use_userInterface "gin-center/internal/domain/interface/user"
	"gin-center/internal/types/constants"
	"gin-center/pkg/security/useJwt"
	"os"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Container 应用程序的依赖注入容器
type Container struct {
	Config        *config.GlobalConfig                   // 应用程序配置
	DB            *gorm.DB                               // 数据库连接
	Redis         *redis.Client                          // Redis客户端
	Logger        *zaplogger.ServiceLogger               // 修改为自定义日志类型  // 日志记录器
	UserService   use_userInterface.UserServiceInterface // 用户服务接口
	AdminService  *AdminService.AdminService             // 管理员服务接口（保持接口名称不变）
	SystemService *systemService.SystemService           // 系统配置服务
	Validator     *validator.Validate                    // 数据验证器
	JWTConfig     *useJwt.JWTConfig                      // JWT配置
	Cache         cache.Cache                            // 缓存接口
	shutdown      sync.Once                              // 确保关闭操作只执行一次
}

// NewContainer 创建并初始化一个新的依赖注入容器
// 按照以下顺序初始化各个组件：
// 1. 日志系统
// 2. JWT认证
// 3. Redis连接
// 4. 数据库连接
// 5. 缓存系统
// 6. 各类服
func NewContainer(cfg *config.GlobalConfig) (*Container, error) {
	if cfg == nil {
		return nil, errors.New("配置对象不能为空")
	}

	// 初始化基础组件
	logger, jwtSecret, redisClient, db, err := initCoreComponents(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化核心组件失败: %w", err)
	}

	// 初始化缓存实例
	cacheInstance := cache.NewRedisCache(redisClient)

	// 配置JWT
	jwtConfig := &useJwt.JWTConfig{
		SecretKey:  jwtSecret,
		Expiration: time.Hour * 72,
	}

	// 初始化仓储层
	adminRepo := admin.NewAdminRepository(db)
	userRepo := user_repo.NewUserRepository(db)

	// 初始化服务层
	services, err := initServices(&serviceConfig{
		DB:           db,
		Cache:        cacheInstance,
		AdminRepo:    adminRepo,
		UserRepo:     userRepo,
		JWTSecret:    jwtSecret,
		GlobalConfig: cfg,
		RedisClient:  redisClient,
		Logger:       logger,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化服务层失败: %w", err)
	}

	// 初始化验证器
	validatorInstance := validator.New()

	return &Container{
		Config:        cfg,
		DB:            db,
		Redis:         redisClient,
		Logger:        logger,
		UserService:   services.UserService,
		AdminService:  services.AdminService,
		SystemService: services.SystemService,
		Validator:     validatorInstance,
		JWTConfig:     jwtConfig,
		Cache:         cacheInstance,
	}, nil
}

func initCoreComponents(cfg *config.GlobalConfig) (*zaplogger.ServiceLogger, string, *redis.Client, *gorm.DB, error) {
	jwtSecret, err := getJWTSecretWithValidation(cfg)
	if err != nil {
		return nil, "", nil, nil, fmt.Errorf("JWT配置验证失败: %w", err)
	}

	db, err := initDatabase(cfg)
	if err != nil {
		return nil, "", nil, nil, fmt.Errorf("数据库初始化失败: %w", err)
	}

	redisClient, err := initRedis(cfg)
	if err != nil {
		return nil, "", nil, nil, fmt.Errorf("Redis连接失败: %w", err)
	}
	return zaplogger.NewServiceLogger(), jwtSecret, redisClient, db, nil
}

func getJWTSecretWithValidation(cfg *config.GlobalConfig) (string, error) {
	secret := os.Getenv("APP_JWT_SECRET")
	if secret == "" {
		secret = cfg.JWT.SecretKey
	}
	if secret == "" {
		secret = constants.DefaultJWTSecret
		return secret, errors.New("使用默认JWT密钥，建议在生产环境配置APP_JWT_SECRET环境变量")
	}
	return secret, nil
}

func initRedis(cfg *config.GlobalConfig) (*redis.Client, error) {
	redisCfg := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:        redisCfg.Password,
		DB:              redisCfg.DB,
		PoolSize:        redisCfg.PoolSize,
		MinIdleConns:    redisCfg.MinIdleConns,
		MaxConnAge:      redisCfg.MaxConnAge,
		IdleTimeout:     redisCfg.IdleTimeout,
		MaxRetries:      redisCfg.MaxRetries,
		MinRetryBackoff: redisCfg.MinRetryDelay,
		MaxRetryBackoff: redisCfg.MaxRetryDelay,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "Redis连接测试失败 host:%s port:%d", redisCfg.Host, redisCfg.Port)
	}
	return client, nil
}

func getJWTSecret(cfg *config.GlobalConfig) string {
	if secret := os.Getenv("APP_JWT_SECRET"); secret != "" {
		return secret
	}
	if cfg.JWT.SecretKey != "" {
		return cfg.JWT.SecretKey
	}
	return constants.DefaultJWTSecret
}

func initRedisClient(redisCfg *config.RedisConfig) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:        redisCfg.Password,
		DB:              redisCfg.DB,
		PoolSize:        redisCfg.PoolSize,
		MinIdleConns:    redisCfg.MinIdleConns,
		MaxConnAge:      redisCfg.MaxConnAge,
		IdleTimeout:     redisCfg.IdleTimeout,
		MaxRetries:      redisCfg.MaxRetries,
		MinRetryBackoff: redisCfg.MinRetryDelay,
		MaxRetryBackoff: redisCfg.MaxRetryDelay,
	}), nil
}

func initDatabase(cfg *config.GlobalConfig) (*gorm.DB, error) {
	return database.InitDatabase(cfg)
}

// serviceConfig 服务初始化配置
type serviceConfig struct {
	DB           *gorm.DB
	Cache        cache.Cache
	AdminRepo    *admin.AdminRepository
	UserRepo     *user_repo.UserRepository
	JWTSecret    string
	GlobalConfig *config.GlobalConfig
	RedisClient  *redis.Client
	Logger       *zaplogger.ServiceLogger // 修改日志类型
}

// ServiceContainer 服务容器，包含所有初始化的服务实例
type ServiceContainer struct {
	UserService   use_userInterface.UserServiceInterface
	AdminService  *AdminService.AdminService
	SystemService *systemService.SystemService
}

// initServices 初始化应用服务
func initServices(cfg *serviceConfig) (*ServiceContainer, error) {
	adminService := AdminService.NewAdminService(cfg.AdminRepo, &useJwt.JWTConfig{SecretKey: cfg.JWTSecret}, cfg.GlobalConfig, cfg.Logger)
	userService := user_service.NewUserService(cfg.UserRepo, cfg.Logger, &useJwt.JWTConfig{SecretKey: cfg.JWTSecret})
	systemService := systemService.NewSystemService(cfg.RedisClient, cfg.Logger)

	return &ServiceContainer{
		UserService:   userService,
		AdminService:  adminService,
		SystemService: systemService,
	}, nil
}

// WithTransaction 在事务中执行数据库操作
func (c *Container) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if c.DB == nil {
		return errors.New("数据库连接未初始化")
	}
	return c.DB.WithContext(ctx).Transaction(fn)
}

// Close 优雅关闭容器中的所有资源
//
// 按以下顺序关闭组件：
// 1. HTTP服务器
// 2. Redis连接
// 3. 数据库连接
// 4. 日志系统
func (c *Container) Close() {
	c.shutdown.Do(func() {
		// 关闭Redis连接
		if c.Redis != nil {
			if err := c.Redis.Close(); err != nil {
				c.Logger.LogError("关闭Redis连接失败", zap.Error(err))
			}
		}

		// 关闭数据库连接
		if c.DB != nil {
			sqlDB, err := c.DB.DB()
			if err != nil {
				c.Logger.LogError("获取SQL.DB实例失败", zap.Error(err))
				return
			}
			if err := sqlDB.Close(); err != nil {
				c.Logger.LogError("关闭数据库连接失败", zap.Error(err))
			}
		}

		// 同步日志缓冲区
		if c.Logger != nil {
			_ = c.Logger.Sync()
		}
	})
}
