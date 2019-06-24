package main

import (
	"context"
	"log"
	"time"
	"math/rand"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"github.com/future-architect/futureot/exporters/opencensus-go-exporter-p8s-pushgateway"
)

// Create measures. The program will record measures for the size of
// processed videos and the number of videos marked as spam.
var (
	videoCount = stats.Int64("example.com/measures/video_count", "number of processed videos", stats.UnitDimensionless)
	videoSize  = stats.Int64("example.com/measures/video_size", "size of processed video", stats.UnitBytes)
)

func main() {
	ctx := context.Background()

	exporter, err := pushgateway.NewExporter(pushgateway.Options{
		GatewayEndpoint: "localhost:9091",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initializing: export to Prometheus pushgateway (localhost:9091)")

	// Create view to see the number of processed videos cumulatively.
	// Create view to see the amount of video processed
	// Subscribe will allow view data to be exported.
	// Once no longer needed, you can unsubscribe from the view.
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

	view.RegisterExporter(exporter)

	// Set reporting period to report data at every second.
	view.SetReportingPeriod(1 * time.Second)

	// Record some data points...
	for {
		log.Println("sending data to pushgateway")
		stats.Record(ctx, videoCount.M(1), videoSize.M(rand.Int63()))
		<-time.After(time.Millisecond * time.Duration(1+rand.Intn(400)))
	}
}
