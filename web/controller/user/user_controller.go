package user_controller

import (
	"errors"
	"fmt"
	"gin-center/internal/types/auth"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	zaplogger "gin-center/infrastructure/zaplogger"
	use_userInterface "gin-center/internal/domain/interface/user"
	type_response "gin-center/internal/types/response"
	use_response "gin-center/pkg/http/response"
	base_controller "gin-center/web/controller"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserController struct {
	base_controller.BaseController
	userService use_userInterface.UserServiceInterface
}

func NewUserController(logger *zaplogger.ServiceLogger, userServiceInterface use_userInterface.UserServiceInterface) *UserController {
	return &UserController{
		BaseController: base_controller.BaseController{Logger: logger}, // 修改为使用 logger 本身
		userService:    userServiceInterface,
	}
}

// @Summary 用户登录
// @Description 处理用户登录请求，验证用户名和密码，返回JWT令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录请求参数"
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "登录成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Failure 401 {object} type_response.BaseResponse "登录失败"
// @Router /api/v1/user/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	c.Logger.LogInfo("User login attempt", zap.String("ip", ctx.ClientIP()))
	var req auth.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Logger.LogError("Login request validation failed", zap.Error(err))
		use_response.BadRequest(ctx, "Invalid login request")
		return
	}
	result, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		c.Logger.LogError("Login failed", zap.String("username", req.Username), zap.Error(err))
		use_response.Unauthorized(ctx, "Login failed: "+err.Error())
		return
	}
	use_response.Authenticated(ctx, result, "")
}

// @Summary 用户注册
// @Description 处理用户注册请求，创建新用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "注册请求参数"
// @Success 200 {object} type_response.BaseResponse "注册成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Failure 500 {object} type_response.BaseResponse "注册失败"
// @Router /api/v1/user/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	c.Logger.LogInfo("User registration attempt", zap.String("ip", ctx.ClientIP()))
	var req auth.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Logger.LogError("Registration request validation failed", zap.Error(err))
		use_response.BadRequest(ctx, "Invalid registration request")
		return
	}
	if err := c.userService.Register(req.Username, req.Password); err != nil {
		c.Logger.LogError("Registration failed",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		use_response.ServerError(ctx, "Registration failed: "+err.Error())
		return
	}
	use_response.Success(ctx, "Registration successful")
}

// validateAvatarFile 验证头像文件的有效性
func (c *UserController) validateAvatarFile(file *multipart.FileHeader) error {
	if file.Size > 5*1024*1024 {
		return errors.New("avatar file size must be less than 5MB")
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	uploadedFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open avatar file: %w", err)
	}
	defer uploadedFile.Close()

	buffer := make([]byte, 512)
	if _, err = uploadedFile.Read(buffer); err != nil {
		return fmt.Errorf("failed to read avatar file: %w", err)
	}

	fileType := http.DetectContentType(buffer)
	if !allowedTypes[fileType] {
		return errors.New("only jpeg, png, and gif images are allowed")
	}

	return nil
}

// @Summary 获取用户个人资料
// @Description 获取当前登录用户的个人资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} type_response.BaseResponse{data=type_response.UserResponse} "获取成功"
// @Failure 404 {object} type_response.BaseResponse "用户不存在"
// @Router /api/v1/user/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	user, err := c.userService.GetUserByID(ctx, userID)
	if err != nil {
		c.Logger.LogError("Failed to get user profile", zap.Uint("user_id", userID), zap.Error(err))
		use_response.NotFound(ctx, "User profile not found")
		return
	}
	profile := type_response.UserResponse{
		ID:        strconv.FormatUint(uint64(user.ID), 10),
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	use_response.Success(ctx, profile)
}

// @Summary 更新用户个人资料
// @Description 更新当前登录用户的个人资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body type_response.UpdateUserProfileRequest true "更新资料请求参数"
// @Success 200 {object} type_response.BaseResponse "更新成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Router /api/v1/user/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req type_response.UpdateUserProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Logger.LogError("Invalid profile update request", zap.Uint("user_id", userID), zap.Error(err))
		use_response.BadRequest(ctx, "Invalid profile update request: "+err.Error())
		return
	}
	err := c.userService.UpdateUserProfile(ctx, userID, &req)
	if err != nil {
		c.Logger.LogError("Failed to update user profile", zap.Uint("user_id", userID), zap.Error(err))
		use_response.BadRequest(ctx, "Profile update failed: "+err.Error())
		return
	}
	use_response.Success(ctx, gin.H{"message": "Profile updated successfully"})
}

// @Summary 获取用户列表
// @Description 分页获取用户列表信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码，默认1" default(1)
// @Param size query int false "每页数量，默认10" default(10)
// @Param username query string false "用户名筛选"
// @Success 200 {object} type_response.BaseResponse{data=type_response.UserListResponse} "获取成功"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/user/list [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	username := ctx.Query("username")
	query := map[string]interface{}{}
	if username != "" {
		query["username"] = username
	}

	c.Logger.LogInfo("Listing users", zap.Int("page", page), zap.Int("size", size), zap.String("username", username))
	users, err := c.userService.ListUsers(ctx, page, size, query)
	if err != nil {
		c.Logger.LogError("Failed to list users", zap.Error(err))
		use_response.ServerError(ctx, "Failed to list users: "+err.Error())
		return
	}
	use_response.Success(ctx, users)
}

// @Summary 上传用户头像
// @Description 上传并更新当前登录用户的头像
// @Tags 用户管理
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param avatar formData file true "用户头像文件（支持jpg、png、gif，小于5MB）"
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "上传成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/user/avatar [post]
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	file, err := ctx.FormFile("avatar")
	if err != nil {
		c.Logger.LogError("Failed to get avatar file", zap.Uint("user_id", userID), zap.Error(err))
		use_response.BadRequest(ctx, "Invalid avatar file")
		return
	}

	if err := c.validateAvatarFile(file); err != nil {
		c.Logger.LogError("Avatar file validation failed", zap.Uint("user_id", userID), zap.Error(err))
		use_response.BadRequest(ctx, err.Error())
		return
	}

	filename := fmt.Sprintf("avatar_%d_%s%s",
		userID,
		uuid.New().String(),
		filepath.Ext(file.Filename),
	)
	avatarPath := filepath.Join("uploads", "avatars", filename)

	if err := os.MkdirAll(filepath.Dir(avatarPath), os.ModePerm); err != nil {
		c.Logger.LogError("Failed to create avatar directory", zap.Uint("user_id", userID), zap.Error(err))
		use_response.ServerError(ctx, "Failed to create avatar directory")
		return
	}

	if err := ctx.SaveUploadedFile(file, avatarPath); err != nil {
		c.Logger.LogError("Failed to save avatar file", zap.Uint("user_id", userID), zap.Error(err))
		use_response.ServerError(ctx, "Failed to save avatar file")
		return
	}

	err = c.userService.UpdateUserAvatar(ctx, userID, "/"+avatarPath)
	if err != nil {
		c.Logger.LogError("Failed to update user avatar in database", zap.Uint("user_id", userID), zap.Error(err))
		use_response.ServerError(ctx, "Failed to update user avatar")
		return
	}

	use_response.Success(ctx, gin.H{
		"avatar_url": "/" + avatarPath,
		"message":    "Avatar uploaded successfully",
	})
}

// @Summary 修改用户密码
// @Description 修改当前登录用户的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body type_response.ChangePasswordRequest true "修改密码请求参数"
// @Success 200 {object} type_response.BaseResponse "修改成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Router /api/v1/user/password [put]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req type_response.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Logger.LogError("Invalid password change request", zap.Uint("user_id", userID), zap.Error(err))
		use_response.BadRequest(ctx, "Invalid password change request: "+err.Error())
		return
	}
	err := c.userService.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.Logger.LogError("Password change failed", zap.Uint("user_id", userID), zap.Error(err))
		use_response.BadRequest(ctx, "Password change failed: "+err.Error())
		return
	}
	use_response.Success(ctx, gin.H{"message": "Password changed successfully"})
}
