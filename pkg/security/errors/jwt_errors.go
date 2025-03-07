package security_errors

import (
	infraredErrors "gin-center/infrastructure/errors"
)

// JWTError 定义JWT相关错误
type JWTError struct {
	Code    infraredErrors.AuthErrorCode
	Message string
}

// Error 实现error接口
func (e *JWTError) Error() string {
	return e.Message
}

// NewJWTError 创建一个新的JWT错误
func NewJWTError(code infraredErrors.AuthErrorCode) *JWTError {
	errMap := map[infraredErrors.AuthErrorCode]string{
		infraredErrors.ErrCodeTokenEmpty:       "认证令牌不能为空",
		infraredErrors.ErrCodeTokenInvalid:     "无效的认证令牌",
		infraredErrors.ErrCodePermissionDenied: "权限不足",
		infraredErrors.ErrCodeTokenExpired:     "令牌已过期",
		infraredErrors.ErrCodeTokenRevoked:     "令牌已被撤销",
	}

	return &JWTError{
		Code:    code,
		Message: errMap[code],
	}
}

// 预定义的JWT错误
var (
	ErrInvalidToken = NewJWTError(infraredErrors.ErrCodeTokenInvalid)
	ErrExpiredToken = NewJWTError(infraredErrors.ErrCodeTokenExpired)
	ErrRevokedToken = NewJWTError(infraredErrors.ErrCodeTokenRevoked)
	ErrEmptyToken   = NewJWTError(infraredErrors.ErrCodeTokenEmpty)
)