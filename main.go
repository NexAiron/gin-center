// Package main 是应用程序的入口点
// 负责初始化应用程序、设置路由并启动HTTP服务器
package main

import (
	"gin-center/infrastructure/bootstrap"
	use_routes "gin-center/web/routes"

	"go.uber.org/zap"
)

// main 函数是应用程序的入口点
// 负责初始化应用、设置路由、启动服务器并在程序结束时进行清理
func main() {
	// 初始化应用程序
	app, err := bootstrap.InitializeApp()
	if err != nil {
		panic(err)
	}
	// 确保在程序结束时进行资源清理
	defer app.Cleanup()

	// 设置路由
	use_routes.SetupRoutes(app.Engine, app.Container)

	// 启动HTTP服务器
	if err := app.Server.Start(app.Engine); err != nil {
		zap.L().Fatal("Server error", zap.Error(err))
	}

	// 优雅关闭服务器
	if err := app.Server.Shutdown(app.Context); err != nil {
		zap.L().Fatal("Server shutdown error", zap.Error(err))
	}
}
