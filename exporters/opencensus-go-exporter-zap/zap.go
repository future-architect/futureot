// zap package provides OpenCensus exporter to zap logging library
//
// It makes the local development easily.
//
// zap information is here: https://godoc.org/go.uber.org/zap
package zap

import (
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

type zapTraceExporter struct{
	logger *zap.Logger
}

// NewZapTraceExporter returns new exporter for opencensus tracing.
// It generates logger by using zap.NewDevelopment().
//
// Use it like this:
//
//   import (
//      "go.opencensus.io/trace"
//   	"github.com/future-architect/futureot/exporters/opencensus-go-exporter-zap"
//   )
//
//   func main() {
//       trace.RegisterExporter(zap.NewZapTraceExporter())
//   }
func NewZapTraceExporter() trace.Exporter {
	logger, _ := zap.NewDevelopment()
	return &zapTraceExporter{
		logger: logger,
	}
}

// NewZapTraceExpoerterWith is similar to NewZapTraceExporter, but you can specify
// well configured logger you want.
func NewZapTraceExpoerterWith(logger *zap.Logger) trace.Exporter {
	return &zapTraceExporter{
		logger: logger,
	}
}

func (z zapTraceExporter) ExportSpan(sd *trace.SpanData) {
	if !sd.IsSampled() {
		return
	}

	fields := []zap.Field{
		zap.String("Name", sd.Name),
		zap.String("TraceID", sd.TraceID.String()),
		zap.String("SpanID", sd.SpanID.String()),
		zap.Time("Start", sd.StartTime),
		zap.Time("End", sd.EndTime),
	}

	for key, att := range sd.Attributes {
		fields = append(fields, zap.Any(key, att))
	}

	z.logger.Info("Span", fields...)
}
