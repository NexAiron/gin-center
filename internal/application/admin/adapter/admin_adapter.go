// Package admin_adapter 实现管理员服务适配器，用于将领域服务适配到控制器层接口
package admin_adapter

import (
	use_AdminInterface "gin-center/internal/domain/interface/admin"
)

// adminServiceAdapter 实现了AdminServiceInterface接口的适配器结构体
type adminServiceAdapter struct {
	adminService use_AdminInterface.AdminServiceInterface
}

// NewAdminServiceAdapter 创建一个新的管理员服务适配器实例
// @Summary 创建管理员服务适配器
// @Description 创建一个新的管理员服务适配器实例
// @Tags Admin
// @Param adminService body use_AdminInterface.AdminServiceInterface true "管理员领域服务接口"
// @Return use_AdminInterface.AdminServiceInterface 实现了管理员服务接口的适配器实例
func NewAdminServiceAdapter(adminService use_AdminInterface.AdminServiceInterface) use_AdminInterface.AdminServiceInterface {
	return &adminServiceAdapter{adminService: adminService}
}

// Login 管理员登录方法
// @Summary 管理员登录
// @Description 处理管理员登录请求
// @Tags Admin
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 401 {object} error "登录失败"
// @Router /admin/login [post]
func (a *adminServiceAdapter) Login(username, password string) (string, map[string]interface{}, error) {
	token, adminInfo, err := a.adminService.Login(username, password)
	if err != nil {
		return "", nil, err
	}
	return token, adminInfo, nil
}

// Register 管理员注册方法
// @Summary 管理员注册
// @Description 注册新管理员账户
// @Tags Admin
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {string} string "注册成功"
// @Failure 400 {object} error "注册失败"
// @Router /admin/register [post]
func (a *adminServiceAdapter) Register(username, password string) error {
	return a.adminService.Register(username, password)
}

// GetAdminInfo 获取管理员信息
// @Summary 获取管理员信息
// @Description 获取指定管理员的详细信息
// @Tags Admin
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} map[string]interface{} "管理员信息"
// @Failure 404 {object} error "未找到管理员"
// @Router /admin/info/{username} [get]
func (a *adminServiceAdapter) GetAdminInfo(username string) (*map[string]any, error) {
	return a.adminService.GetAdminInfo(username)
}

// UpdateAdmin 更新管理员信息
// @Summary 更新管理员信息
// @Description 更新指定管理员的信息
// @Tags Admin
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Param updates body map[string]interface{} true "更新字段"
// @Success 200 {string} string "更新成功"
// @Failure 400 {object} error "更新失败"
// @Router /admin/{username} [put]
func (a *adminServiceAdapter) UpdateAdmin(username string, updates map[string]interface{}) error {
	return a.adminService.UpdateAdmin(username, updates)
}

// PaginateAdmins 分页获取管理员列表
// @Summary 分页获取管理员列表
// @Description 分页获取系统中的管理员列表
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {array} map[string]interface{} "管理员列表"
// @Success 200 {integer} int64 "总数量"
// @Router /admin/list [get]
func (a *adminServiceAdapter) PaginateAdmins(page, pageSize int) ([]map[string]any, int64, error) {
	return a.adminService.PaginateAdmins(page, pageSize)
}
