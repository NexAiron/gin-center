package type_response

import "time"

// BaseResponse 标准API响应结构
// 包含状态码、消息和核心数据
// 扩展字段Token用于需要返回令牌的接口
// RequestID和TimeStamp已移至中间件统一处理
type BaseResponse struct {
	Code    int         `json:"code"`            // 业务状态码
	Message string      `json:"message"`         // 提示信息
	Data    interface{} `json:"data,omitempty"`  // 核心业务数据
	Token   string      `json:"token,omitempty"` // 认证令牌
}

type ListResponse struct {
	Total int64 `json:"total"`
	Page  int   `json:"page" validate:"min=1"`
	Size  int   `json:"size" validate:"min=1"`
	Items any   `json:"items"`
}

type UserResponse struct {
	ID          string     `json:"id" validate:"required"`
	Username    string     `json:"username" validate:"required"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Phone       string     `json:"phone" validate:"omitempty,len=11"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP string     `json:"last_login_ip,omitempty"`
	CreatedAt   time.Time  `json:"created_at" validate:"required"`
	UpdatedAt   time.Time  `json:"updated_at" validate:"required"`
}
type UserListResponse struct {
	ListResponse
	Items []UserResponse `json:"items"`
}
type UpdateUserProfileRequest struct {
	Nickname string `json:"nickname" binding:"required" validate:"required,min=2,max=32"`
}
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" validate:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required" validate:"required,min=6"`
}
