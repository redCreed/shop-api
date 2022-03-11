package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"mxshop-api/order-web/global"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d",global.ServerConfig.JaegerConfig.Host,
					global.ServerConfig.JaegerConfig.Port),
			},

			ServiceName: global.ServerConfig.JaegerConfig.Name,
		}

		trace, closer, err := cfg.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
		if err != nil {
			panic(err)
		}
		defer closer.Close()
		opentracing.InitGlobalTracer(trace)

		startSpan := trace.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()
		c.Set("tracer", trace)
		c.Set("parentSpan", startSpan)
		c.Next()
	}
}