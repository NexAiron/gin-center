package use_tracing

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// InitTracer 初始化Jaeger链路追踪器
// serviceName: 服务名称，用于在Jaeger UI中标识服务
// 返回值：
// - opentracing.Tracer: 全局追踪器实例
// - io.Closer: 用于清理资源的closer
// - error: 初始化过程中的错误信息
func InitTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, // 采样率为100%
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://localhost:14268/api/traces", // Jaeger collector的默认地址
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, nil, fmt.Errorf("初始化Jaeger追踪器失败: %v", err)
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

// Trace 返回一个Gin中间件，用于处理请求的链路追踪
// 该中间件会为每个请求创建一个新的span，并记录请求的关键信息
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取全局追踪器
		tracer := opentracing.GlobalTracer()

		// 生成请求ID并设置到上下文和响应头中
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 尝试从请求头中提取父span上下文
		spanCtx, _ := tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)

		// 创建新的span
		span := tracer.StartSpan(
			c.Request.URL.Path,
			ext.RPCServerOption(spanCtx),
			opentracing.Tag{Key: "request_id", Value: requestID},
		)
		defer span.Finish()

		// 设置span的标签
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		ext.Component.Set(span, "gin")

		// 将span和tracer存储到上下文中
		c.Set("span", span)
		c.Set("tracer", tracer)

		// 继续处理请求
		c.Next()

		// 记录响应状态码
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
	}
}
