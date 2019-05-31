# opencensus-go-exporter-zap

This package provides [OpenCensus](https://opencensus.io/) exporter to [zap](https://godoc.org/go.uber.org/zap) logging library.
It makes the local development easily.

## Sample

```go
package main

import (
	"context"
	
	"go.opencensus.io/trace"
	"github.com/future-architect/futureot/exporters/opencensus-go-exporter-zap"
)

func main () {
	trace.RegisterExporter(zap.NewZapTraceExporter())
	
	ctx, span := trace.StartSpan(context.Background(), "trace-sample")
	defer span.End()
	
	// Do any heavy task
}
```

### Output

```text
2019-05-31T18:12:10.209+0900	INFO	opencensus-go-exporter-zap/zap.go:42	Span
    {"Name": "trace-sample",
     "TraceID": "faaf52233ea266c205db8613ad69128a", "SpanID": "5396d405ba3e3def",
     "Start": "2019-05-31T18:12:09.155+0900", "End": "2019-05-31T18:12:10.209+0900"}
```

## Reference

### ``NewZapTraceExporter()``

``NewZapTraceExporter()`` returns new exporter for opencensus tracing.
It generates logger by using ``zap.NewDevelopment()``.

### ``NewZapTraceExpoerterWith(logger *zap.Logger)``
 
``NewZapTraceExpoerterWith()`` is similar to ``NewZapTraceExporter()``, but you can specify
well configured logger you want.

## License

Apache 2