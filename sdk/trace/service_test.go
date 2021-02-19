package trace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"testing"
)

func TestJager(t *testing.T) {
	tracer, closer := initJaeger("jaeger-demo")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer) //StartspanFromContext创建新span时会用到

	span := tracer.StartSpan("span_root")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	r1 := foo3("Hello foo3", ctx)
	fmt.Println(r1)
	span.Finish()

}
