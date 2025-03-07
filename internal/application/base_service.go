// Package use_Baseservice 提供基础服务层实现
package use_Baseservice

import (
	"context"
	"errors"
	"fmt"
	"gin-center/infrastructure/cache"
	infraErrors "gin-center/infrastructure/errors"
	UserModel "gin-center/internal/domain/model/user"
	"gin-center/internal/types/models/structs"
	security_types "gin-center/pkg/security/types"
	"gin-center/pkg/security/useJwt"
	"gin-center/pkg/utils/validator"
	"time"

	"gin-center/infrastructure/zaplogger"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// BaseService 提供基础服务功能，包括数据库操作、缓存管理、日志记录等
type BaseService struct {
	DB     *gorm.DB                 // 数据库连接实例
	Cache  cache.Cache              // 缓存提供者接口
	Logger *zaplogger.ServiceLogger // 服务日志记录器
	config *BaseServiceConfig       // 服务配置信息
}

// BaseServiceConfig 基础服务配置
type BaseServiceConfig struct {
	DB     *gorm.DB
	Cache  cache.Cache
	Logger *zaplogger.ServiceLogger
}

// NewBaseService 创建BaseService实例
func NewBaseService(config *BaseServiceConfig) *BaseService {
	if config == nil {
		panic("service configuration cannot be nil")
	}
	return &BaseService{
		DB:     config.DB,
		Cache:  config.Cache,
		Logger: zaplogger.NewServiceLogger(),
		config: config,
	}
}

// WithTransaction 在事务中执行数据库操作
func (s *BaseService) WithTransaction(fn func(tx *gorm.DB) error) error {
	return s.DB.Transaction(fn)
}

// WithContext 使用上下文执行数据库操作
func (s *BaseService) WithContext(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// CacheKey 生成缓存键
func (s *BaseService) CacheKey(prefix string, id interface{}) string {
	return fmt.Sprintf("%s:%v", prefix, id)
}

// GetCache 从缓存中获取数据，支持上下文控制和错误处理
// 参数:
//   - ctx: 上下文，如果为nil则使用默认上下文
//   - key: 缓存键
//   - value: 用于存储获取到的缓存值的指针
//
// 返回:
//   - error: 操作过程中的错误，如果未找到数据则返回 ErrNotFound
func (s *BaseService) GetCache(ctx context.Context, key string, value interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}
	result, err := s.Cache.Get(ctx, key)
	if err != nil {
		return s.handleCacheError("获取缓存", key, err)
	}
	if result == nil {
		return infraErrors.ErrNotFound
	}
	// 将缓存结果复制到传入的value中
	if err := s.Cache.Unmarshal(result, value); err != nil {
		return s.handleCacheError("解析缓存数据", key, err)
	}
	return nil
}

// SetCache 设置缓存数据
func (s *BaseService) SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := s.Cache.Set(ctx, key, value, expiration); err != nil {
		return s.handleCacheError("设置缓存", key, err)
	}
	return nil
}

// DeleteCache 删除缓存数据
func (s *BaseService) DeleteCache(ctx context.Context, key string) error {
	if err := s.Cache.Delete(ctx, key); err != nil {
		return s.handleCacheError("删除缓存", key, err)
	}
	return nil
}

// handleCacheError 处理缓存操作错误
func (s *BaseService) handleCacheError(operation, key string, err error) error {
	s.Logger.LogError(fmt.Sprintf("缓存操作失败: %s, key: %s", operation, key), zap.Skip(), zap.Error(err))
	return infraErrors.ErrCache
}

// ValidateUserInput 验证用户输入
func (s *BaseService) ValidateUserInput(username, password string) error {
	if err := validator.ValidateUsername(username); err != nil {
		return err
	}
	return validator.ValidatePassword(password)
}

// HashPassword 使用bcrypt对密码进行哈希处理
func (s *BaseService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}
	return string(hash), nil
}

// ComparePassword 比较密码是否匹配
func (s *BaseService) ComparePassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("密码不匹配: %w", err)
	}
	return nil
}

// GenerateToken 生成JWT令牌
// 参数:
//   - jwtConfig: JWT配置信息
//   - userID: 用户ID
//   - username: 用户名
//
// 返回:
//   - string: 生成的JWT令牌
//   - error: 生成过程中的错误
func (s *BaseService) GenerateToken(jwtConfig *useJwt.JWTConfig, userID uint, username string) (string, error) {
	// 验证 JWT 配置
	if jwtConfig == nil || jwtConfig.SecretKey == "" {
		s.Logger.LogError("JWT配置无效", zap.String("username", username), zap.Error(errors.New("JWT配置或密钥为空")))
		return "", fmt.Errorf("JWT配置无效")
	}

	// 验证输入参数
	if username == "" {
		s.Logger.LogError("用户名不能为空", zap.Skip())
		return "", fmt.Errorf("用户名不能为空")
	}

	// 创建令牌声明
	claims := &structs.UserClaims{
		UserID: fmt.Sprintf("%d", userID),
	}

	// 生成令牌
	token, err := jwtConfig.GenerateTokenWithClaims(claims)
	if err != nil {
		s.Logger.LogError("生成JWT令牌失败", zap.String("username", username), zap.Error(err))
		return "", fmt.Errorf("生成JWT令牌失败: %w", err)
	}

	s.Logger.LogInfo("JWT令牌生成成功", zap.String("username", username))
	return token, nil
}

// ValidateToken 验证JWT令牌
// 参数:
//   - jwtConfig: JWT配置信息
//   - tokenString: 待验证的JWT令牌字符串
//
// 返回:
//   - *structs.Claims: 解析出的令牌声明信息
//   - error: 验证过程中的错误
func (s *BaseService) ValidateToken(jwtConfig *useJwt.JWTConfig, tokenString string) (*security_types.UserClaims, error) {
	if jwtConfig == nil {
		return nil, fmt.Errorf("JWT配置不能为空")
	}
	claims, err := jwtConfig.ParseToken(tokenString)
	if err != nil {
		s.Logger.LogError("验证JWT令牌失败", zap.String("tokenString", tokenString), zap.Error(err))
		return nil, fmt.Errorf("验证JWT令牌失败: %w", err)
	}
	return claims, nil
}

// CheckUserExists 检查用户是否已存在
func (s *BaseService) CheckUserExists(findUserFunc func(ctx context.Context, username string) (*UserModel.User, error), username string) error {
	if findUserFunc == nil {
		return fmt.Errorf("findUserFunc cannot be nil")
	}
	existingUser, err := findUserFunc(context.Background(), username)
	if err != nil {
		s.Logger.LogError("检查用户存在性失败", zap.String("username", username), zap.Error(err))
		return fmt.Errorf("检查用户存在性失败: %w", err)
	}
	if existingUser != nil {
		return &infraErrors.AppError{Code: 409, Message: "用户已存在"}
	}
	return nil
}

// HandleError 统一处理服务层错误
func (s *BaseService) HandleError(err error) error {
	if err != nil {
		s.Logger.LogError("Service error", zap.Skip(), zap.Error(err))
		return &infraErrors.AppError{Code: 500, Message: err.Error()}
	}
	return nil
}
