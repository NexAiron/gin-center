// Package headers 提供统一的HTTP请求头和响应头处理
package use_headers

import (
	"github.com/gin-gonic/gin"
)

// 常用的HTTP头常量
const (
	Authorization   = "Authorization"
	ContentType     = "Content-Type"
	Accept          = "Accept"
	XRequestID      = "X-Request-ID"
	XTraceID        = "X-Trace-ID"
	BearerPrefix    = "Bearer "
	ApplicationJSON = "application/json"
	ApplicationForm = "application/x-www-form-urlencoded"
)

// GetAuthorizationToken 从请求头中获取认证令牌
func GetAuthorizationToken(c *gin.Context) string {
	bearer := c.GetHeader(Authorization)
	if len(bearer) > len(BearerPrefix) && bearer[:len(BearerPrefix)] == BearerPrefix {
		return bearer[len(BearerPrefix):]
	}
	return ""
}

// SetTraceHeaders 设置追踪相关的响应头
func SetTraceHeaders(c *gin.Context, requestID, traceID string) {
	if requestID != "" {
		c.Header(XRequestID, requestID)
	}
	if traceID != "" {
		c.Header(XTraceID, traceID)
	}
}

// GetRequestHeaders 获取请求头信息
func GetRequestHeaders(c *gin.Context) map[string]string {
	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	return headers
}
