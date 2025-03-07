// 自定义错误类型、错误处理函数和错误码定义
package infraErrors

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// AuthErrorCode 定义认证相关错误类型
type AuthErrorCode int

const (
	ErrCodeTokenEmpty       AuthErrorCode = 4001
	ErrCodeTokenInvalid     AuthErrorCode = 4002
	ErrCodePermissionDenied AuthErrorCode = 4003
	ErrCodeTokenExpired     AuthErrorCode = 4004
	ErrCodeTokenRevoked     AuthErrorCode = 4005
)

// AppError 定义了应用程序的自定义错误类型
// 包含错误码、错误信息、详细信息和原始错误
type AppError struct {
	// Code 表示错误码
	Code int
	// Message 表示错误信息
	Message string
	// Details 包含额外的错误详情
	Details interface{}
	// Err 保存原始错误
	Err error
}

// Error 实现error接口
func (e *AppError) Error() string {
	return e.Message
}

// NewError 创建一个新的AppError实例
// code: 错误码
// message: 错误信息
func NewError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// ErrorResponse 定义了API错误响应的结构
type ErrorResponse struct {
	// Code 表示错误码
	Code int `json:"code"`
	// Message 表示错误信息
	Message string `json:"message"`
	// TraceID 用于追踪请求
	TraceID string `json:"trace_id"`
	// Details 包含错误的详细信息（可选）
	Details interface{} `json:"details,omitempty"`
	// DateTime 表示错误发生的时间
	DateTime string `json:"datetime"`
}

// 预定义的错误常量
var (
	// ErrNotFound 表示请求的资源未找到
	ErrNotFound = NewError(404, "资源未找到")
	// ErrValidation 表示请求参数验证失败
	ErrValidation = NewError(400, "参数验证失败")
	// ErrUnauthorized 表示未经授权的访问
	ErrUnauthorized = NewError(401, "未经授权的访问")
	// ErrForbidden 表示禁止访问
	ErrForbidden = NewError(403, "禁止访问")
	// ErrBadRequest 表示错误的请求
	ErrBadRequest = NewError(400, "错误的请求")
	// ErrTooManyRequests 表示请求过于频繁
	ErrTooManyRequests = NewError(429, "请求过于频繁")
	// ErrInvalidCredentials 表示无效的凭证
	ErrInvalidCredentials = NewError(401, "无效的凭证")
	// ErrDatabase 表示数据库操作错误
	ErrDatabase = NewError(500, "数据库操作错误")
	// ErrCache 表示缓存操作错误
	ErrCache = NewError(500, "缓存操作错误")
	// ErrInternalServer 表示服务器内部错误
	ErrInternalServer = NewError(500, "服务器内部错误")
	// ErrServiceUnavailable 表示服务不可用
	ErrServiceUnavailable = NewError(503, "服务不可用")
	// ErrTimeout 表示请求超时
	ErrTimeout = NewError(504, "请求超时")
	// ErrRateLimit 表示请求达到限制
	ErrRateLimit = NewError(429, "请求达到限制")
	// ErrCircuitBreaker 表示熔断器触发
	ErrCircuitBreaker = NewError(503, "服务熔断")
	// ErrResourceConflict 表示资源冲突
	ErrResourceConflict = NewError(600, "资源冲突")
	// ErrResourceExhausted 表示资源耗尽
	ErrResourceExhausted = NewError(601, "资源耗尽")
	// ErrBusinessValidation 表示业务验证失败
	ErrBusinessValidation = NewError(602, "业务验证失败")
	// ErrDataInconsistency 表示数据不一致
	ErrDataInconsistency = NewError(603, "数据不一致")

	// 新增JWT相关错误
	ErrTokenInvalid = NewError(401, "token无效")
	ErrTokenExpired = NewError(401, "token过期")
	ErrTokenRevoked = NewError(401, "token已撤销")
	ErrTokenEmpty   = NewError(400, "token不能为空")

	// 新增认证相关错误
	ErrUserNotFound     = NewError(404, "用户不存在")
	ErrUsernameExists   = NewError(400, "用户名已存在")
	ErrPasswordMismatch = NewError(401, "密码错误")
)

// ErrorHandler 返回一个Gin中间件，用于全局错误处理
// 该中间件可以捕获panic，记录错误日志，并返回统一的错误响应
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetString("trace_id")
		if traceID == "" {
			traceID = uuid.New().String()
			c.Set("trace_id", traceID)
			c.Header("X-Trace-ID", traceID)
		}
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				logFields := logrus.Fields{
					"trace_id":    traceID,
					"error":       err,
					"stack_trace": stack,
					"path":        c.Request.URL.Path,
					"method":      c.Request.Method,
					"client_ip":   c.ClientIP(),
					"user_agent":  c.Request.UserAgent(),
					"referer":     c.Request.Referer(),
					"timestamp":   time.Now().Format(time.RFC3339),
				}
				// 记录请求头信息，排除敏感信息
				headers := make(map[string]string)
				for k, v := range c.Request.Header {
					if k != "Authorization" && k != "Cookie" {
						headers[k] = v[0]
					}
				}
				logFields["headers"] = headers
				logFields["query_params"] = c.Request.URL.Query()
				logrus.WithFields(logFields).Error("Panic recovered")

				// 处理不同类型的错误
				var appError *AppError
				switch e := err.(type) {
				case *AppError:
					appError = e
				case error:
					switch {
					case e.Error() == "record not found":
						appError = ErrNotFound
					case e.Error() == "validation failed":
						appError = ErrValidation
					default:
						appError = &AppError{
							Code:    ErrInternalServer.Code,
							Message: ErrInternalServer.Message,
							Err:     e,
						}
					}
				default:
					appError = &AppError{
						Code:    ErrInternalServer.Code,
						Message: ErrInternalServer.Message,
						Err:     fmt.Errorf("%v", err),
					}
				}

				// 返回统一的错误响应
				c.JSON(appError.Code, ErrorResponse{
					Code:     appError.Code,
					Message:  appError.Message,
					TraceID:  traceID,
					Details:  appError.Details,
					DateTime: time.Now().Format(time.RFC3339),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

const (
	ErrCodeInvalidUsername          AuthErrorCode = 1001
	ErrCodeUsernameValidationFailed AuthErrorCode = 1002
	ErrCodeInvalidPassword          AuthErrorCode = 1003
)
