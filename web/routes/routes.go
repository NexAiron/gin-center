// Package routes 提供应用程序的路由配置功能
// 包括API路由注册、中间件应用和Swagger文档配置
package use_routes

import (
	"gin-center/infrastructure/container"
	"gin-center/infrastructure/zaplogger"
	admin_controller "gin-center/web/controller/admin"
	system_controller "gin-center/web/controller/system"
	user_controller "gin-center/web/controller/user"
	use_AuthMiddleware "gin-center/web/middleware/auth"

	"gin-center/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes 配置应用程序的所有路由
func SetupRoutes(r *gin.Engine, container *container.Container) {
	zapLogger := zaplogger.NewServiceLogger()

	// 初始化所有控制器
	userCtrl := user_controller.NewUserController(zapLogger, container.UserService)
	adminCtrl := admin_controller.NewAdminController(container.AdminService, zapLogger)
	systemCtrl := system_controller.NewSystemController(container.SystemService, zapLogger)

	// 基础路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger文档配置
	docs.SwaggerInfo.Title = "Gin-Center API"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 路由组
	apiV1 := r.Group("/api/v1")
	{
		// 认证相关路由
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", userCtrl.Login)
			authGroup.POST("/register", userCtrl.Register)
		}

		// 管理员专属路由
		adminGroup := apiV1.Group("/admin")
		adminGroup.Use(use_AuthMiddleware.AdminAuth(zapLogger))
		{
			adminGroup.GET("/users", adminCtrl.PaginateAdmins)
			adminGroup.GET("/profile", adminCtrl.GetAdminInfo)
			adminGroup.PUT("/profile", adminCtrl.UpdateAdmin)
		}

		// 需要JWT认证的通用路由
		authRequired := apiV1.Group("")
		authRequired.Use(use_AuthMiddleware.JWTAuth(container.Config, zapLogger))
		{
			// 用户个人中心
			userCenter := authRequired.Group("/user")
			{
				userCenter.GET("/profile", userCtrl.GetProfile)
				userCenter.PUT("/profile", userCtrl.UpdateProfile)
				userCenter.POST("/avatar", userCtrl.UploadAvatar)
			}

			// 系统管理接口
			systemGroup := authRequired.Group("/system")
			{
				systemGroup.GET("/config", systemCtrl.GetSystemConfig)
				systemGroup.PUT("/config", systemCtrl.UpdateSystemConfig)
				systemGroup.GET("/metrics", systemCtrl.GetSystemMetrics)
			}
		}
	}
}
