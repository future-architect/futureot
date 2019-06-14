package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/future-architect/futureot/occonfig"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	videoCount = stats.Int64("example.com/measures/video_count", "number of processed videos", stats.UnitDimensionless)
	videoSize  = stats.Int64("example.com/measures/video_size", "size of processed video", stats.UnitBytes)
)

func main() {
	ctx := context.Background()
	occonfig.UseFlag(occonfig.Stats)

	flag.Parse()

	oc, err := occonfig.Init(occonfig.Stats)
	if err != nil {
		log.Fatalf("Fail to initialize OpenCensus: %v", err)
	}
	defer oc.Close()

	if err = view.Register(
		&view.View{
			Name:        "video_count",
			Description: "number of videos processed over time",
			Measure:     videoCount,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        "video_size",
			Description: "processed video size over time",
			Measure:     videoSize,
			Aggregation: view.Distribution(0, 1<<16, 1<<32),
		},
	); err != nil {
		log.Fatalf("Cannot register the view: %v", err)
	}

	view.SetReportingPeriod(1 * time.Second)

	for {
		stats.Record(ctx, videoCount.M(1), videoSize.M(rand.Int63()))
		<-time.After(time.Millisecond * time.Duration(1+rand.Intn(400)))
	}
}
