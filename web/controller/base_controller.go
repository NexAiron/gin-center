package base_controller

import (
	"errors"
	"net/http"
	"strconv"

	"gin-center/internal/types/constants"
	use_response "gin-center/pkg/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	zaplogger "gin-center/infrastructure/zaplogger"
	type_response "gin-center/internal/types/response"
)

// BaseController 基础控制器结构体，提供通用的控制器功能
type BaseController struct {
	Logger *zaplogger.ServiceLogger // 使用封装后的日志记录器
}

// ParsePaginationParams 解析分页参数
func (c *BaseController) ParsePaginationParams(ctx *gin.Context) (int, int, error) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 || pageSize < 1 || pageSize > 100 {
		return 0, 0, errors.New("invalid pagination parameters")
	}
	return page, pageSize, nil
}

// NewBaseController 创建一个新的基础控制器实例
func NewBaseController(logger *zaplogger.ServiceLogger) *BaseController {
	return &BaseController{
		Logger: logger,
	}
}

// SendResponse 发送统一格式的HTTP响应
func (c *BaseController) SendResponse(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, type_response.BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// HandleError 统一错误处理方法
func (c *BaseController) HandleError(ctx *gin.Context, err error) {
	c.Logger.LogError("请求处理发生错误",
		zap.Error(err),
		zap.String("request_id", ctx.GetString("request_id")),
	)
	c.SendResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
}

// SendSuccess 发送成功响应
func (c *BaseController) SendSuccess(ctx *gin.Context, data interface{}) {
	c.SendResponse(ctx, http.StatusOK, "success", data)
}

// SendBadRequest 发送请求参数错误响应
func (c *BaseController) SendBadRequest(ctx *gin.Context, message string) {
	c.SendResponse(ctx, http.StatusBadRequest, message, nil)
}

// SendUnauthorized 发送未授权响应
func (c *BaseController) SendUnauthorized(ctx *gin.Context, message string) {
	c.SendResponse(ctx, http.StatusUnauthorized, message, nil)
}

// SendForbidden 发送禁止访问响应
func (c *BaseController) SendForbidden(ctx *gin.Context, message string) {
	c.SendResponse(ctx, http.StatusForbidden, message, nil)
}

// SendConflict 发送资源冲突响应
func (c *BaseController) SendConflict(ctx *gin.Context, message string) {
	c.SendResponse(ctx, http.StatusConflict, message, nil)
}

// SendNotFound 发送资源未找到响应
func (c *BaseController) SendNotFound(ctx *gin.Context, message string) {
	c.SendResponse(ctx, http.StatusNotFound, message, nil)
}

// HandleLogin 通用登录处理方法
func (c *BaseController) HandleLogin(ctx *gin.Context, authService interface {
	Login(username, password string) (string, map[string]interface{}, error)
}) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.SendBadRequest(ctx, "无效的请求参数")
		return
	}

	token, data, err := authService.Login(req.Username, req.Password)
	if err != nil {
		c.SendUnauthorized(ctx, "认证失败")
		return
	}

	c.SendSuccess(ctx, gin.H{
		"token": token,
		"data":  data,
		"user": gin.H{
			"username": req.Username,
		},
	})
}

// 从上下文获取当前用户名
func (c *BaseController) GetCurrentUsername(ctx *gin.Context) (string, error) {
	username, exists := ctx.Get("username")
	if !exists {
		c.Logger.LogError("用户身份信息缺失", zap.Error(errors.New("context中未找到用户信息")))
		use_response.Unauthorized(ctx, "用户未认证")
		return "", constants.ErrUnauthorized
	}
	return username.(string), nil
}

// 请求参数验证方法
func (c *BaseController) ValidateRequest(ctx *gin.Context, req interface{}) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		c.Logger.LogError("请求参数绑定失败", zap.Error(err))
		use_response.BadRequest(ctx, "无效的请求参数")
		return err
	}
	return nil
}
