// 用户类型  通用状态
package enums

// UserType 定义用户类型
// 可能值：admin(管理员), regular(普通用户), guest(访客)
type UserType string

const (
	// UserTypeAdmin 系统管理员，拥有最高权限
	UserTypeAdmin UserType = "admin"

	// UserTypeRegular 普通注册用户，基础权限
	UserTypeRegular UserType = "regular"

	// UserTypeGuest 未登录访客，受限权限
	UserTypeGuest UserType = "guest"
)

// IsValid 验证UserType值是否合法
// 返回true表示是有效的用户类型
func (t UserType) IsValid() bool {
	switch t {
	case UserTypeAdmin, UserTypeRegular, UserTypeGuest:
		return true
	default:
		return false
	}
}

// Status 表示通用状态类型
// 用于系统各种业务状态管理
type Status string

const (
	// StatusActive 活跃/生效状态
	StatusActive Status = "active"

	// StatusInactive 未激活/失效状态
	StatusInactive Status = "inactive"

	// StatusPending 等待处理状态
	StatusPending Status = "pending"

	// StatusProcessing 正在处理中
	StatusProcessing Status = "processing"

	// StatusCompleted 处理完成
	StatusCompleted Status = "completed"

	// StatusCancelled 已取消
	StatusCancelled Status = "cancelled"

	// StatusFailed 处理失败
	StatusFailed Status = "failed"
)

// IsValid 验证Status值是否合法
// 返回true表示是有效的状态值
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusPending,
		StatusProcessing, StatusCompleted, StatusCancelled, StatusFailed:
		return true
	default:
		return false
	}
}
