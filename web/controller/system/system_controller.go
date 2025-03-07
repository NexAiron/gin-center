package system_controller

import (
	zaplogger "gin-center/infrastructure/zaplogger"
	systemService "gin-center/internal/application/system/system_service"
	"gin-center/internal/types/system"
	base_controller "gin-center/web/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemController struct {
	*base_controller.BaseController
	systemService *systemService.SystemService
}

func NewSystemController(systemService *systemService.SystemService, logger *zaplogger.ServiceLogger) *SystemController {
	return &SystemController{
		BaseController: base_controller.NewBaseController(logger),
		systemService:  systemService,
	}
}
// @Summary 获取系统信息
// @Description 获取系统基本信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "获取成功"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/system/info [get]
func (c *SystemController) GetSystemInfo(ctx *gin.Context) {
	data, err := c.systemService.GetSystemInfo()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}
	c.SendResponse(ctx, http.StatusOK, "success", data)
}
// @Summary 获取系统配置
// @Description 获取系统配置信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} type_response.BaseResponse{data=system.SystemConfig} "获取成功"
// @Router /api/v1/system/config [get]
func (c *SystemController) GetSystemConfig(ctx *gin.Context) {
	data := c.systemService.GetSystemConfig()
	c.SendResponse(ctx, http.StatusOK, "success", data)
}
// @Summary 更新系统配置
// @Description 更新系统配置信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param config body system.SystemConfig true "系统配置信息"
// @Success 200 {object} type_response.BaseResponse "更新成功"
// @Failure 400 {object} type_response.BaseResponse "请求参数错误"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/system/config [put]
func (c *SystemController) UpdateSystemConfig(ctx *gin.Context) {
	var config system.SystemConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		c.SendResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}
	if err := c.validateSystemConfig(&config); err != nil {
		c.SendResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if err := c.systemService.UpdateSystemConfig(config); err != nil {
		c.HandleError(ctx, err)
		return
	}
	c.SendResponse(ctx, http.StatusOK, "success", nil)
}
// @Summary 获取系统指标
// @Description 获取系统运行指标信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "获取成功"
// @Failure 500 {object} type_response.BaseResponse "服务器错误"
// @Router /api/v1/system/metrics [get]
func (c *SystemController) GetSystemMetrics(ctx *gin.Context) {
	data, err := c.systemService.GetSystemMetrics()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}
	c.SendResponse(ctx, http.StatusOK, "success", data)
}
// @Summary 获取系统健康状态
// @Description 获取系统健康检查信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} type_response.BaseResponse{data=map[string]interface{}} "获取成功"
// @Router /api/v1/system/health [get]
func (c *SystemController) GetSystemHealth(ctx *gin.Context) {
	data := c.systemService.GetSystemHealth()
	c.SendResponse(ctx, http.StatusOK, "success", data)
}
func (c *SystemController) validateSystemConfig(config *system.SystemConfig) error {
	return nil
}
