package auth

import (
	infraErrors "gin-center/infrastructure/errors"
	"gin-center/internal/types/request"
)

type AuthError struct {
	Code    infraErrors.AuthErrorCode
	Message string
}

var (
	ErrTokenInvalid = NewAuthError(infraErrors.ErrCodeTokenInvalid)
	ErrTokenRevoked = NewAuthError(infraErrors.ErrCodeTokenRevoked)
)

func (e *AuthError) Error() string {
	return e.Message
}
func NewAuthError(code infraErrors.AuthErrorCode) *AuthError {
	errMap := map[infraErrors.AuthErrorCode]string{
		infraErrors.ErrCodeTokenEmpty:       "认证令牌不能为空",
		infraErrors.ErrCodeTokenInvalid:     "无效的认证令牌",
		infraErrors.ErrCodePermissionDenied: "权限不足",
		infraErrors.ErrCodeTokenExpired:     "令牌已过期",
		infraErrors.ErrCodeTokenRevoked:     "令牌已被撤销",
	}

	return &AuthError{
		Code:    code,
		Message: errMap[code],
	}
}

type RegisterRequest struct {
	request.BaseRequest
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	request.BaseRequest
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
