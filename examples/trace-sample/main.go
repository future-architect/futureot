package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/future-architect/futureot/occonfig"
	"go.opencensus.io/trace"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"time"
)

func main() {
	occonfig.UseKingpin(occonfig.Trace)

	kingpin.Parse()

	oc, err := occonfig.Init(occonfig.Trace)
	if err != nil {
		log.Fatalf("Fail to initialize OpenCensus: %v", err)
	}
	defer oc.Close()

	ctx, span := trace.StartSpan(context.Background(), "trace-sample")
	defer span.End()

	for i := 0; i < 10; i++ {
		doWork(ctx)
	}
}

func doWork(ctx context.Context) {
	// 4. Start a child span. This will be a child span because we've passed
	// the parent span's ctx.
	_, span := trace.StartSpan(ctx, "doWork")
	// 5a. Make the span close at the end of this function.
	defer span.End()

	fmt.Println("doing busy work")
	time.Sleep(80 * time.Millisecond)
	buf := bytes.NewBuffer([]byte{0xFF, 0x00, 0x00, 0x00})
	num, err := binary.ReadVarint(buf)
	if err != nil {
		// 6. Set status upon error
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
	}

	// 7. Annotate our span to capture metadata about our operation
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("bytes to int", num),
	}, "Invoking doWork")
	time.Sleep(20 * time.Millisecond)
}
