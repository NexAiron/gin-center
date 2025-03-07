// Package http 提供HTTP服务器的核心功能实现
package use_http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	use_config "gin-center/configs/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HTTPServer 定义HTTP服务器的接口
type HTTPServer interface {
	// Start 启动HTTP服务器
	Start(router *gin.Engine) error
	// Shutdown 优雅关闭HTTP服务器
	Shutdown(ctx context.Context) error
}

// httpServer 实现HTTPServer接口的服务器结构体
type httpServer struct {
	config *use_config.ServerConfig
	server *http.Server
	logger *zap.Logger
}

// NewHTTPServer 创建一个新的HTTP服务器实例
func NewHTTPServer(
	port int,
	readTimeout, writeTimeout, idleTimeout, readHeaderTimeout, shutdownTimeout time.Duration,
	maxHeaderBytes, maxRequestsPerConn int,
) HTTPServer {
	cfg := &use_config.ServerConfig{
		Port:               port,
		ReadTimeout:        readTimeout,
		WriteTimeout:       writeTimeout,
		IdleTimeout:        idleTimeout,
		ReadHeaderTimeout:  readHeaderTimeout,
		ShutdownTimeout:    shutdownTimeout,
		MaxHeaderBytes:     maxHeaderBytes,
		MaxRequestsPerConn: maxRequestsPerConn,
	}

	return &httpServer{
		config: cfg,
		logger: zap.L(),
	}
}

// Start 实现HTTPServer接口的Start方法
func (s *httpServer) Start(router *gin.Engine) error {
	// 配置HTTP服务器
	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.config.Port),
		Handler:           router,
		ReadTimeout:       s.config.ReadTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		MaxHeaderBytes:    s.config.MaxHeaderBytes,
	}

	// 异步启动服务器
	go s.startServer()
	return nil
}

// startServer 启动HTTP服务器并处理错误
func (s *httpServer) startServer() {
	s.logger.Info("服务器正在启动", zap.Int("端口", s.config.Port))

	// 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待关闭信号
	sig := <-sigChan
	s.logger.Info("收到关闭信号", zap.String("信号", sig.String()))

	// 执行关闭操作
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("服务器关闭失败", zap.Error(err))
	}
}

// Shutdown 实现HTTPServer接口的Shutdown方法
func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// shutdownServer 执行服务器关闭操作
func (s *httpServer) shutdownServer(ctx context.Context) error {
	s.logger.Info("正在关闭服务器...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("服务器关闭失败", zap.Error(err))
		return fmt.Errorf("服务器关闭失败: %w", err)
	}
	return nil
}
