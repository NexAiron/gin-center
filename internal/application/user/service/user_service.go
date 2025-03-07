package UserService

import (
	"context"
	"errors"
	"fmt"
	user_repo "gin-center/infrastructure/repository/user"
	use_Baseservice "gin-center/internal/application"
	use_userInterface "gin-center/internal/domain/interface/user"
	UserModel "gin-center/internal/domain/model/user"
	"gin-center/internal/types/models/structs"
	type_response "gin-center/internal/types/response"
	useJwt "gin-center/pkg/security/useJwt"
	"strconv"

	"gin-center/infrastructure/zaplogger"

	"go.uber.org/zap"
)

// UserService 实现用户服务接口
type UserService struct {
	baseService *use_Baseservice.BaseService
	userRepo    *user_repo.UserRepository
	logger      *zaplogger.ServiceLogger
	jwtConfig   *useJwt.JWTConfig
}

// NewUserService 创建新的用户服务实例
func NewUserService(userRepo *user_repo.UserRepository, logger *zaplogger.ServiceLogger, jwtConfig *useJwt.JWTConfig) use_userInterface.UserServiceInterface {
	return &UserService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
		logger:    logger,
	}
}

// Register 用户注册
func (s *UserService) Register(username, password string, extraFields ...interface{}) error {
	s.logger.LogInfo("Processing user registration", zap.String("username", username))

	if err := s.baseService.ValidateUserInput(username, password); err != nil {
		s.logger.LogWarn("Invalid registration input", zap.String("username", username), zap.Error(err))
		return err
	}

	existingUser, err := s.userRepo.FindByUsername(context.Background(), username)
	if err == nil && existingUser != nil {
		s.logger.LogWarn("Username already exists", zap.String("username", username))
		return errors.New("username already exists")
	}

	hashedPassword, err := s.baseService.HashPassword(password)
	if err != nil {
		s.logger.LogError("Failed to hash password", zap.Error(err))
		return err
	}

	user := &UserModel.User{
		Username: username,
		Password: hashedPassword,
	}

	return s.userRepo.Register(context.Background(), user)
}

// Login 用户登录
func (s *UserService) Login(username, password string) (map[string]interface{}, error) {
	s.logger.LogInfo("User login attempt", zap.String("username", username))

	if err := s.baseService.ValidateUserInput(username, password); err != nil {
		s.logger.LogWarn("Invalid login input", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	user, err := s.userRepo.FindByUsername(context.Background(), username)
	if err != nil {
		s.logger.LogWarn("User not found", zap.String("username", username), zap.Error(err))
		return nil, errors.New("invalid username or password")
	}

	if err := s.baseService.ComparePassword(password, user.Password); err != nil {
		s.logger.LogWarn("Invalid password attempt", zap.String("username", username))
		return nil, errors.New("invalid username or password")
	}
	userClaims := &structs.UserClaims{
		UserID: strconv.FormatUint(uint64(user.ID), 10),
	}
	tokenString, err := s.jwtConfig.GenerateTokenWithClaims(userClaims)
	if err != nil {
		s.logger.LogError("Failed to generate token", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	s.logger.LogInfo("User login successful", zap.String("username", username))
	response := map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
	}

	return response, nil
}

// ValidateToken 验证JWT令牌
func (s *UserService) ValidateToken(tokenString string) (*structs.UserClaims, error) {
	securityClaims, err := s.jwtConfig.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 将security_types.UserClaims转换为structs.UserClaims
	return &structs.UserClaims{
		BaseClaims: structs.BaseClaims{
			ID:       securityClaims.BaseClaims.ID,
			Username: securityClaims.BaseClaims.Username,
		},
		UserID: securityClaims.UserID,
	}, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*UserModel.User, error) {
	s.logger.LogDebug("Getting user by ID", zap.Uint("user_id", id))
	return s.userRepo.FindByID(ctx, id)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, user *UserModel.User) error {
	s.logger.LogInfo("Updating user information", zap.Uint("user_id", user.ID))
	_, err := s.userRepo.Update(ctx, user)
	return err
}

// ListUsers 分页获取用户列表
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int, query map[string]interface{}) (*type_response.UserListResponse, error) {
	s.logger.LogInfo("Listing users", zap.Int("page", page), zap.Int("page_size", pageSize))

	if page < 1 || pageSize < 1 {
		err := errors.New("invalid pagination parameters")
		s.logger.LogWarn("Invalid pagination parameters", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return nil, err
	}

	users, total, err := s.userRepo.ListUsers(ctx, page, pageSize, query)
	if err != nil {
		s.logger.LogError("Failed to list users", zap.Error(err))
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	if users == nil {
		users = make([]*UserModel.User, 0)
	}

	userResponses := make([]type_response.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = type_response.UserResponse{
			ID:       strconv.FormatUint(uint64(u.ID), 10),
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
		}
	}

	return &type_response.UserListResponse{
		ListResponse: type_response.ListResponse{
			Total: int64(total),
			Page:  page,
			Size:  pageSize,
		},
		Items: userResponses,
	}, nil
}

// UpdateUserAvatar 更新用户头像
func (s *UserService) UpdateUserAvatar(ctx context.Context, userID uint, avatarPath string) error {
	s.logger.LogInfo("Updating user avatar", zap.Uint("user_id", userID), zap.String("avatar_path", avatarPath))
	return s.userRepo.UpdateAvatar(ctx, userID, avatarPath)
}

// UpdateUserProfile 更新用户个人资料
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uint, profile *type_response.UpdateUserProfileRequest) error {
	s.logger.LogInfo("Updating user profile", zap.Uint("user_id", userID))

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.LogError("Failed to find user for profile update", zap.Uint("user_id", userID), zap.Error(err))
		return err
	}

	// 更新用户资料字段
	if profile.Nickname != "" {
		user.Nickname = profile.Nickname
	}

	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.LogError("Failed to update user profile", zap.Uint("user_id", userID), zap.Error(err))
	}
	return err
}

// ChangePassword 修改用户密码
func (s *UserService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	s.logger.LogInfo("Changing user password", zap.Uint("user_id", userID))

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.LogError("Failed to find user for password change", zap.Uint("user_id", userID), zap.Error(err))
		return err
	}

	if err := s.baseService.ComparePassword(oldPassword, user.Password); err != nil {
		s.logger.LogWarn("Invalid old password", zap.Uint("user_id", userID))
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := s.baseService.HashPassword(newPassword)
	if err != nil {
		s.logger.LogError("Failed to hash new password", zap.Uint("user_id", userID), zap.Error(err))
		return err
	}

	user.Password = hashedPassword
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.LogError("Failed to update password", zap.Uint("user_id", userID), zap.Error(err))
	}
	return err
}
