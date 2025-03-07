// response类型 通用状态
package constants

import (
	"errors"
	"net/http"
)

var (
	ErrUserExists         = errors.New("用户已存在")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrInvalidCredentials = errors.New("无效的凭证")
	ErrUserInactive       = errors.New("账号已禁用")
	ErrUnauthorized       = errors.New("未授权的访问")
)

const DefaultJWTSecret = "gin-center-default-secret"

type ResponseCode int

const (
	CodeSuccess         ResponseCode = 200
	CodeCreated         ResponseCode = 201
	CodeNoContent       ResponseCode = 204
	CodeBadRequest      ResponseCode = 400
	CodeUnauthorized    ResponseCode = 401
	CodeForbidden       ResponseCode = 403
	CodeNotFound        ResponseCode = 404
	CodeConflict        ResponseCode = 409
	CodeTooManyRequests ResponseCode = 429
	CodeServerError     ResponseCode = 500
)

type MessageConstants string

const (
	MsgSuccess       MessageConstants = "success"
	MsgCreated       MessageConstants = "created"
	MsgAuthenticated MessageConstants = "authenticated"
	MsgNoContent     MessageConstants = "no content"
	MsgBadRequest    MessageConstants = "invalid request"
)

type ResponseCategory int

const (
	CategorySuccess     ResponseCategory = 2
	CategoryClientError ResponseCategory = 4
	CategoryServerError ResponseCategory = 5
)

var (
	statusMap = map[int]int{
		int(CodeSuccess):         http.StatusOK,
		int(CodeCreated):         http.StatusCreated,
		int(CodeNoContent):       http.StatusNoContent,
		int(CodeBadRequest):      http.StatusBadRequest,
		int(CodeUnauthorized):    http.StatusUnauthorized,
		int(CodeForbidden):       http.StatusForbidden,
		int(CodeNotFound):        http.StatusNotFound,
		int(CodeConflict):        http.StatusConflict,
		int(CodeTooManyRequests): http.StatusTooManyRequests,
		int(CodeServerError):     http.StatusInternalServerError,
	}
)

func MapResponseCodeToHTTPStatus(code int) int {
	if status, exists := statusMap[code]; exists {
		return status
	}
	return http.StatusOK
}
func ToResponseCode(code int) ResponseCode {
	return ResponseCode(code)
}
