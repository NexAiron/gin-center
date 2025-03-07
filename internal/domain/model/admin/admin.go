// Package admin_model 定义管理员领域模型
package admin_model

import (
	baseModel "gin-center/internal/domain/model/base"
	"time"
)

// Admin 管理员模型
type Admin struct {
	baseModel.BaseModel
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	LastLoginAt time.Time `json:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip"`
	IsAdmin     int        `json:"is_admin"`
}

// NewAdmin 创建新的管理员
func NewAdmin(username string, password string) *Admin {
	return &Admin{
		Username: username,
		Password: password,
		Avatar:   "default_avatar.png",
		BaseModel: baseModel.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// TableName 返回数据库表名
func (a Admin) TableName() string {
	return "sys_users"
}
