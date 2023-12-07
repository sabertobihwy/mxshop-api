package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	jeagerClient "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"go.uber.org/zap"
)

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		zap.S().Infof("jaeper tracing....1")
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jeagerClient.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: "192.168.2.112:6831",
			},
			ServiceName: "mxshop-api", // tracer Name
		}

		jLogger := jaegerlog.StdLogger

		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jLogger))
		if err != nil {
			panic(err)
			return
		}
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()
		// a span: invoking an interface

		span := opentracing.StartSpan(c.Request.URL.Path)
		defer span.Finish()
		c.Set("tracer", tracer)
		c.Set("parentSpan", span)
		c.Next()
	}
}
