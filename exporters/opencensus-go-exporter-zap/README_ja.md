# opencensus-go-exporter-zap

このパッケージは、[zap](https://godoc.org/go.uber.org/zap)のロガーに出力する、[OpenCensus](https://opencensus.io/)エクスポーターを実装しています。

## サンプル

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
	
	// 何か重い仕事をする
}
```

### 出力

```text
2019-05-31T18:12:10.209+0900	INFO	opencensus-go-exporter-zap/zap.go:42	Span
    {"Name": "trace-sample",
     "TraceID": "faaf52233ea266c205db8613ad69128a", "SpanID": "5396d405ba3e3def",
     "Start": "2019-05-31T18:12:09.155+0900", "End": "2019-05-31T18:12:10.209+0900"}
```

## リファレンス

### ``NewZapTraceExporter()``

``NewZapTraceExporter()``はOpenCensusのトレーシングのエクスポーターを返します。
この関数は、出力先のロガーを``zap.NewDevelopment()``で作成します。

### ``NewZapTraceExpoerterWith(logger *zap.Logger)``
 
``NewZapTraceExpoerterWith()``は``NewZapTraceExporter()``と似ていますが、内部でロガーは作成せずに、
外部で指定したロガーを出力先に設定したエクスポーターを返します。

## License

Apache 2