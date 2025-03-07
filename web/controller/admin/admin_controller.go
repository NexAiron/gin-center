package admin_controller

import (
	"errors"
	"gin-center/infrastructure/zaplogger"
	use_AdminInterface "gin-center/internal/domain/interface/admin"
	"gin-center/internal/types/auth"
	"gin-center/internal/types/constants"
	use_response "gin-center/pkg/http/response"
	base_controller "gin-center/web/controller"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AdminController 管理员控制器，处理管理员相关的HTTP请求
type AdminController struct {
	base_controller.BaseController
	adminService use_AdminInterface.AdminServiceInterface
}

// NewAdminController 创建新的管理员控制器实例
func NewAdminController(adminServiceInterface use_AdminInterface.AdminServiceInterface, logger *zaplogger.ServiceLogger) *AdminController {
	return &AdminController{
		BaseController: *base_controller.NewBaseController(logger),
		adminService:   adminServiceInterface,
	}
}
// @Summary 管理员登录
// @Description 处理管理员登录请求，验证用户名和密码，返回JWT令牌
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录请求参数"
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "登录成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Failure 401 {object} type_response.BaseResponse "登录失败"
// @Router /api/v1/admin/login [post]
func (c *AdminController) Login(ctx *gin.Context) {
	c.Logger.LogDebug("管理员登录尝试", zap.String("ip", ctx.ClientIP()))
	c.HandleLogin(ctx, c.adminService)
}
// @Summary 管理员注册
// @Description 处理管理员注册请求，创建新的管理员账户
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "注册请求参数"
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "注册成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Failure 409 {object} type_response.BaseResponse "用户已存在"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/admin/register [post]
func (c *AdminController) Register(ctx *gin.Context) {
	c.Logger.LogDebug("管理员注册尝试", zap.String("ip", ctx.ClientIP()))
	var req auth.RegisterRequest
	if err := c.BaseController.ValidateRequest(ctx, &req); err != nil {
		c.Logger.LogError("请求参数验证失败", zap.String("username", req.Username), zap.Error(err))
		return
	}

	// 执行注册
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger.LogError("密码加密失败", zap.String("username", req.Username), zap.Error(err))
		use_response.ServerError(ctx, "系统错误")
		return
	}
	if err := c.adminService.Register(req.Username, string(hashedPassword)); err != nil {
		c.Logger.LogError("注册操作失败", zap.String("username", req.Username), zap.Error(err))
		if errors.Is(err, constants.ErrUserExists) {
			c.SendConflict(ctx, "用户已存在")
		} else {
			use_response.ServerError(ctx, "注册处理失败")
		}
		return
	}
	use_response.Success(ctx, map[string]interface{}{
		"username": req.Username,
		"message":  "注册成功",
	})
}
// @Summary 获取管理员信息
// @Description 获取当前登录管理员的详细信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "获取成功"
// @Failure 401 {object} type_response.BaseResponse "未授权"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/admin/info [get]
func (c *AdminController) GetAdminInfo(ctx *gin.Context) {
	usernameStr, err := c.BaseController.GetCurrentUsername(ctx)
	if err != nil {
		return
	}
	c.Logger.LogDebug("获取管理员信息", zap.String("username", usernameStr))

	// 获取管理员信息
	adminInfo, err := c.adminService.GetAdminInfo(usernameStr)
	if err != nil {
		c.Logger.LogError("获取管理员信息失败", zap.String("username", usernameStr), zap.Error(err))
		use_response.ServerError(ctx, "获取管理员信息失败："+err.Error())
		return
	}

	use_response.Success(ctx, adminInfo)
}
// @Summary 更新管理员信息
// @Description 更新当前登录管理员的信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body map[string]interface{} true "更新信息参数"
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "更新成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Failure 401 {object} type_response.BaseResponse "未授权"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/admin/update [put]
func (c *AdminController) UpdateAdmin(ctx *gin.Context) {
	usernameStr, err := c.BaseController.GetCurrentUsername(ctx)
	if err != nil {
		return
	}
	c.Logger.LogDebug("更新管理员信息请求", zap.String("username", usernameStr))

	// 解析更新数据
	var updates map[string]interface{}
	if err := c.BaseController.ValidateRequest(ctx, &updates); err != nil {
		return
	}

	// 执行更新
	if err := c.adminService.UpdateAdmin(usernameStr, updates); err != nil {
		c.Logger.LogError("更新失败", zap.String("username", usernameStr), zap.Error(err))
		use_response.BadRequest(ctx, "更新失败："+err.Error())
		return
	}

	use_response.Success(ctx, map[string]interface{}{
		"username": usernameStr,
		"message":  "更新成功",
	})
}
// @Summary 获取管理员列表
// @Description 分页获取管理员列表信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页数量，默认10" default(10)
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "获取成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/admin/list [get]
func (c *AdminController) PaginateAdmins(ctx *gin.Context) {
	// 获取分页参数
	page, pageSize, err := c.BaseController.ParsePaginationParams(ctx)
	if err != nil {
		c.Logger.LogError("分页参数解析失败", zap.Skip(), zap.Error(err))
		use_response.BadRequest(ctx, "无效的分页参数")
		return
	}

	c.Logger.LogDebug("分页获取管理员列表", zap.Int("page", page), zap.Int("page_size", pageSize))

	// 查询管理员列表
	admins, total, err := c.adminService.PaginateAdmins(page, pageSize)
	if err != nil {
		c.Logger.LogError("获取管理员列表失败", zap.Skip(), zap.Error(err))
		use_response.ServerError(ctx, "获取管理员列表失败："+err.Error())
		return
	}

	c.Logger.LogInfo("获取管理员列表成功", zap.Int64("total", total), zap.Int("count", len(admins)))
	use_response.Success(ctx, map[string]interface{}{
		"total":     total,
		"items":     admins,
		"page":      page,
		"page_size": pageSize,
	})
}
