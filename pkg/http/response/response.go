// Package use_response 提供统一的HTTP响应处理功能
package use_response

import (
	type_response "gin-center/internal/types/response"

	"github.com/gin-gonic/gin"
)

// Option 定义响应选项的函数类型
type Option func(*type_response.BaseResponse)

// WithToken 设置响应中的token
func WithToken(token string) Option {
	return func(r *type_response.BaseResponse) {
		r.Token = token
	}
}

// New 创建新的响应对象
func New(code int, message string, data interface{}, opts ...Option) *type_response.BaseResponse {
	resp := &type_response.BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
	for _, opt := range opts {
		opt(resp)
	}
	return resp
}

// Send 发送HTTP响应
func Send(c *gin.Context, resp *type_response.BaseResponse) {
	c.JSON(resp.Code, resp)
}
func ResponseHandler(c *gin.Context, code int, message string, data interface{}, opts ...Option) {
	resp := New(code, message, data, opts...)
	Send(c, resp)
}
func Success(c *gin.Context, data interface{}, opts ...Option) {
	ResponseHandler(c, 200, "success", data, opts...)
}
func Created(c *gin.Context, data interface{}, opts ...Option) {
	ResponseHandler(c, 201, "created", data, opts...)
}
func Forbidden(c *gin.Context, message string, opts ...Option) {
	ResponseHandler(c, 403, message, nil, opts...)
}
func NotFound(c *gin.Context, message string, opts ...Option) {
	ResponseHandler(c, 404, message, nil, opts...)
}
func ServiceUnavailable(c *gin.Context, message string, opts ...Option) {
	ResponseHandler(c, 503, message, nil, opts...)
}
func Authenticated(c *gin.Context, data map[string]any, token string) {
	ResponseHandler(c, 200, "authenticated", data, WithToken(token))
}
func Unauthorized(c *gin.Context, message string) {
	ResponseHandler(c, 401, message, nil)
}
func BadRequest(c *gin.Context, message string) {
	ResponseHandler(c, 400, message, nil)
}
func ServerError(c *gin.Context, message string) {
	ResponseHandler(c, 500, message, nil)
}
