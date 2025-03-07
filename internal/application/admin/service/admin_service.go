// Package admin_service 实现管理员服务
package admin_service

import (
	"context"
	"fmt"
	"gin-center/configs/config"
	"gin-center/infrastructure/repository/admin"
	zaplogger "gin-center/infrastructure/zaplogger"
	use_Baseservice "gin-center/internal/application"
	AdminModel "gin-center/internal/domain/model/admin"
	"gin-center/internal/types/constants"
	"gin-center/internal/types/models/structs"
	useJwt "gin-center/pkg/security/useJwt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService 管理员服务结构体，提供管理员相关的核心业务功能
type AdminService struct {
	baseService *use_Baseservice.BaseService
	logger      *zaplogger.ServiceLogger
	adminRepo   *admin.AdminRepository
	jwtConfig   *useJwt.JWTConfig
	config      *config.GlobalConfig
}

// NewAdminService 创建新的管理员服务实例
func NewAdminService(adminRepo *admin.AdminRepository, jwtConfig *useJwt.JWTConfig, config *config.GlobalConfig, logger *zaplogger.ServiceLogger) *AdminService {
	return &AdminService{
		baseService: use_Baseservice.NewBaseService(&use_Baseservice.BaseServiceConfig{}),
		logger:      zaplogger.NewServiceLogger(),
		adminRepo:   adminRepo,
		jwtConfig:   jwtConfig,
		config:      config,
	}
}

// handleLoginError 统一处理登录错误
func (s *AdminService) handleLoginError(err error, logMsg string, errType error, username string) error {
	if err != nil {
		s.logger.LogError(logMsg,
			zap.String("module", "login"),
			zap.String("username", username),
			zap.Error(err))
		return fmt.Errorf("%w: %v", errType, err)
	}
	return nil
}

// validatePassword 验证密码有效性
func (s *AdminService) validatePassword(hashedPassword, inputPassword, username string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err != nil {
		s.logger.LogWarn("密码验证失败",
			zap.String("module", "login"),
			zap.String("username", username),
			zap.Error(err))
		return fmt.Errorf("%w: %v", constants.ErrInvalidCredentials, err)
	}
	return nil
}

// checkAdminStatus 检查管理员状态
func (s *AdminService) checkAdminStatus(admin *AdminModel.Admin, username string) error {
	if admin.Status == 0 {
		s.logger.LogWarn("账号已禁用",
			zap.String("module", "login"),
			zap.String("username", username))
		return fmt.Errorf("%w", constants.ErrUserInactive)
	}
	return nil
}

// Register 管理员注册
// 参数:
//   - username: 用户名
//   - password: 密码
//
// 返回:
//   - error: 注册过程中的错误信息
//
// 更新后的Register方法
func (s *AdminService) Register(username, password string) error {
	return s.withTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := s.baseService.ValidateUserInput(username, password); err != nil {
			return s.handleError(err, "register", username, "输入验证失败")
		}

		hashedPassword, err := s.validateAndHashPassword(password)
		if err != nil {
			return err
		}

		return tx.Create(&AdminModel.Admin{
			Username: username,
			Password: hashedPassword,
		}).Error
	})
}

// 更新后的PaginateAdmins方法

// 删除独立的checkAdminStatus方法（功能已整合到handleError）
func (s *AdminService) handleError(err error, module string, username string, msg string) error {
	if err != nil {
		s.logger.LogError(msg,
			zap.String("module", module),
			zap.String("username", username),
			zap.Error(err))
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

// 合并后的密码处理流程
func (s *AdminService) validateAndHashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), s.handleError(err, "security", "", "密码处理失败")
}

// 用户信息构造模板
func (s *AdminService) buildAdminInfo(admin AdminModel.Admin) map[string]interface{} {
	return map[string]interface{}{
		"id":            admin.ID,
		"username":      admin.Username,
		"nickname":      admin.Nickname,
		"avatar":        admin.Avatar,
		"status":        admin.Status,
		"created_at":    admin.CreatedAt,
		"last_login_at": admin.LastLoginAt,
	}
}

// 事务处理模板
func (s *AdminService) withTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return s.adminRepo.DB.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (s *AdminService) Login(username, password string) (string, map[string]interface{}, error) {
	admin, err := s.adminRepo.FindByUsername(context.Background(), username)
	if err != nil {
		s.logger.LogError("管理员登录失败：用户不存在", zap.String("username", username), zap.Error(err))
		return "", nil, constants.ErrUserNotFound
	}
	if err := s.handleError(err, "login", username, "用户查询失败"); err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", nil, s.handleError(err, "login", username, "密码验证失败")
	}

	if admin.Status == 0 {
		return "", nil, s.handleError(nil, "login", username, "账号已禁用")
	}

	var token string
	err = s.withTransaction(context.Background(), func(tx *gorm.DB) error {
		admin.LastLoginAt = time.Now()
		if err := tx.Save(admin).Error; err != nil {
			return err
		}
		token, err = s.GenerateToken(admin)
		return err
	})

	return token, s.buildAdminInfo(*admin), s.handleError(err, "login", username, "登录流程异常")
}

// UpdateAdmin 更新管理员信息
// 参数:
//   - username: 用户名
//   - updates: 需要更新的字段map，支持更新password、nickname和avatar
//
// 返回:
//   - error: 更新过程中的错误信息
func (s *AdminService) UpdateAdmin(username string, updates map[string]interface{}) error {
	ctx := context.Background()
	admin, err := s.adminRepo.FindByUsername(ctx, username)
	if err != nil {
		s.logger.LogError("查询用户失败", zap.String("username", username), zap.Error(err))
		return fmt.Errorf("查询用户失败: %w", err)
	}
	if password, ok := updates["password"].(string); ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.LogError("密码加密失败", zap.String("username", username), zap.Error(err))
			return fmt.Errorf("密码加密失败: %w", err)
		}
		admin.Password = string(hashedPassword)
	}
	if nickname, ok := updates["nickname"].(string); ok {
		admin.Nickname = nickname
	}
	if avatar, ok := updates["avatar"].(string); ok {
		admin.Avatar = avatar
	}
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		s.logger.LogError("更新用户信息失败", zap.String("username", username), zap.Error(err))
		return fmt.Errorf("更新用户信息失败: %w", err)
	}
	s.logger.LogInfo("更新用户信息成功", zap.String("username", username))
	return nil
}

// GetAdminInfo 获取管理员信息
// 参数:
//   - username: 用户名
//
// 返回:
//   - *map[string]interface{}: 包含管理员详细信息的map
//   - error: 获取过程中的错误信息
func (s *AdminService) GetAdminInfo(username string) (*map[string]interface{}, error) {
	ctx := context.Background()
	admin, err := s.adminRepo.FindByUsername(ctx, username)
	if err != nil {
		s.logger.LogError("查询用户失败", zap.String("username", username), zap.Error(err))
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	// 直接构造并返回管理员信息map
	return &map[string]interface{}{
		"id":            admin.ID,
		"username":      admin.Username,
		"nickname":      admin.Nickname,
		"avatar":        admin.Avatar,
		"status":        admin.Status,
		"is_admin":      admin.IsAdmin,
		"created_at":    admin.CreatedAt,
		"updated_at":    admin.UpdatedAt,
		"last_login_at": admin.LastLoginAt,
		"last_login_ip": admin.LastLoginIP,
	}, nil
}

// PaginateAdmins 分页获取管理员列表
// 参数:
//   - page: 页码，从1开始
//   - pageSize: 每页数量
//
// 返回:
//   - []map[string]interface{}: 管理员列表
//   - int64: 总记录数
//   - error: 获取过程中的错误信息
func (s *AdminService) PaginateAdmins(page, pageSize int) ([]map[string]interface{}, int64, error) {
	ctx := context.Background()
	admins, total, err := s.adminRepo.PaginateAdmins(ctx, page, pageSize)
	if err != nil {
		s.logger.LogError("获取管理员列表失败", zap.Error(err))
		return nil, 0, fmt.Errorf("获取管理员列表失败: %w", err)
	}

	// 构造管理员列表
	result := make([]map[string]interface{}, len(admins))
	for i, admin := range admins {
		result[i] = map[string]interface{}{
			"id":            admin.ID,
			"username":      admin.Username,
			"nickname":      admin.Nickname,
			"avatar":        admin.Avatar,
			"status":        admin.Status,
			"is_admin":      admin.IsAdmin,
			"created_at":    admin.CreatedAt,
			"updated_at":    admin.UpdatedAt,
			"last_login_at": admin.LastLoginAt,
			"last_login_ip": admin.LastLoginIP,
		}
	}
	return result, total, nil
}

// GenerateToken 生成JWT令牌
// 参数:
//   - admin: 管理员实体
//
// 返回:
//   - string: JWT令牌字符串
//   - error: 生成过程中的错误信息
func (s *AdminService) GenerateToken(admin *AdminModel.Admin) (string, error) {
	claims := &structs.AdminClaims{
		BaseClaims: structs.BaseClaims{
			Username: admin.Username,
		},
		IsAdmin: admin.IsAdmin == 1,
	}
	tokenString, err := s.jwtConfig.GenerateTokenWithClaims(claims)
	if err != nil {
		s.logger.LogError("生成令牌失败", zap.String("username", admin.Username), zap.Error(err))
		return "", fmt.Errorf("生成令牌失败: %w", err)
	}
	return tokenString, nil
}
