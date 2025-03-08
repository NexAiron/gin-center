package bootstrap

import (
	"context"
	"fmt"
	"gin-center/configs/config"
	"gin-center/infrastructure/container"
	zaplogger "gin-center/infrastructure/zaplogger"
	use_http "gin-center/pkg/http/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// App 应用程序的核心结构体，包含所有主要组件
type App struct {
	Context   context.Context      // 应用程序的上下文，用于控制生命周期
	Config    *config.GlobalConfig // 应用程序配置
	Container *container.Container // 依赖注入容器
	Server    use_http.HTTPServer  // HTTP服务器实例
	Engine    *gin.Engine          // Gin引擎实例
	cancel    context.CancelFunc   // 用于取消上下文的函数
	once      sync.Once            // 确保Close方法只执行一次的同步原语
}

// InitializeApp 初始化并返回一个新的应用程序实例
func InitializeApp() (*App, error) {
	// 创建带取消功能的上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 加载应用配置
	cfg, err := loadConfig()
	if err != nil {
		cancel() // 发生错误时取消上下文
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 初始化日志系统
	if err := initLogger(cfg); err != nil {
		cancel()
		return nil, fmt.Errorf("初始化日志系统失败: %w", err)
	}

	// 初始化依赖注入容器
	c, err := container.NewContainer(cfg)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("初始化容器失败: %w", err)
	}

	// 初始化Gin引擎
	engine := gin.Default()

	// 配置HTTP服务器
	httpServer := use_http.NewHTTPServer(
		cfg.Server.Port,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		cfg.Server.IdleTimeout,
		cfg.Server.ReadHeaderTimeout,
		cfg.Server.ShutdownTimeout,
		cfg.Server.MaxHeaderBytes,
		cfg.Server.MaxRequestsPerConn,
	)

	return &App{
		Context:   ctx,
		Config:    cfg,
		Container: c,
		Server:    httpServer,
		Engine:    engine,
		cancel:    cancel,
	}, nil
}

// loadConfig 加载应用程序配置
// 从指定的配置目录加载配置文件，并返回解析后的配置对象
func loadConfig() (*config.GlobalConfig, error) {
	cfg, err := config.LoadConfig("configs")
	if err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %w", err)
	}
	return cfg, nil
}

// Close 优雅关闭应用程序
// 确保所有组件按照正确的顺序关闭，避免资源泄漏
func (a *App) Close() {
	a.once.Do(func() {
		// 取消应用程序上下文
		a.cancel()

		// 关闭HTTP服务器
		if a.Server != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = a.Server.Shutdown(ctx)
		}

		// 关闭依赖注入容器
		if a.Container != nil {
			a.Container.Close()
		}
	})
}

// Cleanup 执行最终的清理工作
// 在应用程序退出前调用，确保所有资源都被正确释放
func (a *App) Cleanup() {
	a.Close()
	_ = zap.L().Sync()
}

// initLogger 初始化日志系统
// 配置并初始化zap日志记录器，设置日志级别、输出路径等
func initLogger(cfg *config.GlobalConfig) error {
	// 创建日志目录
	logDir := filepath.Dir(cfg.Log.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return errors.Wrap(err, "创建日志目录失败")
	}

	// 配置日志输出
	logConfig := zap.NewProductionConfig()
	logConfig.OutputPaths = []string{cfg.Log.Filename}
	logConfig.ErrorOutputPaths = []string{cfg.Log.Filename}

	// 设置日志级别
	logLevel := zap.NewAtomicLevel()
	if err := logLevel.UnmarshalText([]byte(strings.ToLower(cfg.Log.Level))); err != nil {
		logLevel.SetLevel(zap.InfoLevel)
	}
	logConfig.Level = logLevel
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 构建日志记录器
	logger, err := logConfig.Build()
	if err != nil {
		return errors.Wrap(zaplogger.ErrInitLoggerFailed, err.Error())
	}

	// 设置全局日志记录器
	zap.ReplaceGlobals(logger)
	return nil
}
