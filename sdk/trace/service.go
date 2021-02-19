package trace

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

type Service struct {
	conf   Config
	tracer opentracing.Tracer
}

func GetService(conf Config) *Service {
	return &Service{conf: conf}
}

func (d *Service) Init() {
	cfg := &config.Configuration{
		ServiceName: d.conf.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%s", d.conf.AgentHost, d.conf.AgentPort),
		},
	}
	var err error
	d.tracer, _, err = cfg.NewTracer()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(d.tracer)
}

func (d *Service) Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		if d.conf.Enable == true {
			var span opentracing.Span

			spCtx, _ := d.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			if spCtx == nil {
				span = d.tracer.StartSpan(c.Request.URL.RequestURI())
			} else {
				span = d.tracer.StartSpan(c.Request.URL.RequestURI(), opentracing.ChildOf(spCtx))
			}
			_ = d.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			span.SetTag("method", c.Request.Method)
			c.Next()
			span.Finish()
		}
	}

}

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func foo3(req string, ctx context.Context) (reply string) {
	//1.创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "span_foo3")
	defer func() {
		//4.接口调用完，在tag中设置request和reply
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	println(req)
	//2.模拟处理耗时
	time.Sleep(time.Second / 2)
	//3.返回reply
	reply = "foo3Reply"
	return
}
