package use_tracer

import (
	use_tracing "gin-center/web/middleware/tracing"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	globalTracer opentracing.Tracer
	globalCloser io.Closer
)

func InitTracer(serviceName string) error {
	if serviceName == "" {
		return errors.New("")
	}
	tracer, closer, err := use_tracing.InitTracer(serviceName)
	if err != nil {
		zap.L().Error("",
			zap.String("service", serviceName),
			zap.Error(err))
		return errors.Wrap(err, "")
	}
	globalTracer = tracer
	globalCloser = closer
	zap.L().Info("", zap.String("service", serviceName))
	return nil
}
func CloseTracer() error {
	if globalCloser != nil {
		if err := globalCloser.Close(); err != nil {
			zap.L().Error("", zap.Error(err))
			return errors.Wrap(err, "")
		}
		zap.L().Info("")
	}
	return nil
}
func GetTracer() opentracing.Tracer {
	return globalTracer
}
