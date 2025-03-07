package UserAdapter

import (
	use_userInterface "gin-center/internal/domain/interface/user"
	UserModel "gin-center/internal/domain/model/user"
	type_response "gin-center/internal/types/response"

	"github.com/gin-gonic/gin"
)

type userServiceAdapter struct {
	userService use_userInterface.UserServiceInterface
}

// NewUserServiceAdapter 创建用户服务适配器
// @Summary 创建用户服务适配器
// @Description 创建一个新的用户服务适配器实例
// @Tags User
// @Param userService body use_userInterface.UserServiceInterface true "用户服务接口"
// @Return use_userInterface.UserServiceInterface 用户服务适配器实例
func NewUserServiceAdapter(userService use_userInterface.UserServiceInterface) *userServiceAdapter {
	return &userServiceAdapter{userService: userService}
}

// Login 用户登录
// @Summary 用户登录
// @Description 处理用户登录请求
// @Tags User
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 401 {object} error "登录失败"
// @Router /user/login [post]
func (a *userServiceAdapter) Login(username, password string) (map[string]interface{}, error) {
	return a.userService.Login(username, password)
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户账户
// @Tags User
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {string} string "注册成功"
// @Failure 400 {object} error "注册失败"
// @Router /user/register [post]
func (a *userServiceAdapter) Register(username, password string) error {
	return a.userService.Register(username, password)
}

// GetUserByID 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户ID获取用户信息
// @Tags User
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} UserModel.User "用户信息"
// @Failure 404 {object} error "用户不存在"
// @Router /user/{id} [get]
func (a *userServiceAdapter) GetUserByID(ctx *gin.Context, userID uint) (*UserModel.User, error) {
	user, err := a.userService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUserProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新用户个人资料
// @Tags User
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Param profile body type_response.UpdateUserProfileRequest true "用户资料"
// @Success 200 {string} string "更新成功"
// @Failure 400 {object} error "更新失败"
// @Router /user/{id}/profile [put]
func (a *userServiceAdapter) UpdateUserProfile(ctx *gin.Context, userID uint, profile *type_response.UpdateUserProfileRequest) error {
	return a.userService.UpdateUserProfile(ctx.Request.Context(), userID, profile)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改用户密码
// @Tags User
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Param old_password body string true "旧密码"
// @Param new_password body string true "新密码"
// @Success 200 {string} string "修改成功"
// @Failure 400 {object} error "修改失败"
// @Router /user/{id}/password [put]
func (a *userServiceAdapter) ChangePassword(ctx *gin.Context, userID uint, oldPassword, newPassword string) error {
	return a.userService.ChangePassword(ctx.Request.Context(), userID, oldPassword, newPassword)
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags User
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param query query object false "查询条件"
// @Success 200 {object} type_response.UserListResponse "用户列表"
// @Router /user/list [get]
func (a *userServiceAdapter) ListUsers(ctx *gin.Context, page, pageSize int, query map[string]interface{}) (*type_response.UserListResponse, error) {
	return a.userService.ListUsers(ctx.Request.Context(), page, pageSize, query)
}

// UpdateUserAvatar 更新用户头像
// @Summary 更新用户头像
// @Description 更新用户头像
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param id path uint true "用户ID"
// @Param avatar formData file true "头像文件"
// @Success 200 {string} string "更新成功"
// @Failure 400 {object} error "更新失败"
// @Router /user/{id}/avatar [put]
func (a *userServiceAdapter) UpdateUserAvatar(ctx *gin.Context, userID uint, avatarPath string) error {
	return a.userService.UpdateUserAvatar(ctx.Request.Context(), userID, avatarPath)
}
