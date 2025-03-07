package use_AuthMiddleware

import (
	"gin-center/configs/config"
	use_headers "gin-center/pkg/http/headers"
	use_response "gin-center/pkg/http/response"
	useJwt "gin-center/pkg/security/useJwt"
	"strings"

	"gin-center/infrastructure/zaplogger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTAuth 统一的JWT认证中间件
func JWTAuth(cfg *config.GlobalConfig, logger *zaplogger.ServiceLogger) gin.HandlerFunc {
	// 初始化JWT配置
	jwtConfig := useJwt.NewJWTConfig(&useJwt.JWTConfig{
		SecretKey:            cfg.JWT.SecretKey,
		AccessTokenLifetime:  cfg.JWT.AccessTokenLifetime,
		RefreshTokenLifetime: cfg.JWT.RefreshTokenLifetime,
		Issuer:               cfg.JWT.Issuer,
		BlacklistCleanupTick: cfg.JWT.BlacklistCleanupTick,
	})
	// 返回中间件处理函数
	return func(c *gin.Context) {
		// 获取并验证Authorization头
		token := use_headers.GetAuthorizationToken(c)
		if token == "" {
			logger.LogWarn("缺少认证头")
			use_response.Unauthorized(c, "缺少认证头")
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.LogWarn("无效的认证头格式")
			use_response.Unauthorized(c, "无效的认证头格式")
			c.Abort()
			return
		}

		// 验证JWT token
		claims, err := jwtConfig.ParseToken(parts[1])
		if err != nil {
			logger.LogWarn("Token验证失败", zap.Error(err))
			use_response.Unauthorized(c, "无效的Token")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.ID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// AdminAuth 管理员权限验证中间件
func AdminAuth(logger *zaplogger.ServiceLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("username")
		if !exists {
			use_response.Unauthorized(c, "未经授权的访问")
			c.Abort()
			return
		}

		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			use_response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// SuperAuth 超级管理员权限验证中间件
func SuperAuth(cfg *config.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			use_response.Unauthorized(c, "未经授权的访问")
			c.Abort()
			return
		}

		if username.(string) != cfg.App.SuperAdmin {
			use_response.Forbidden(c, "需要超级管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
