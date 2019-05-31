module stats_sample

go 1.12

replace (
    github.com/future-architect/futureot/occonfig => ../../occonfig
    github.com/future-architect/futureot/exporters/opencensus-go-exporter-zap => ../../exporters/opencensus-go-exporter-zap
)

require (
	go.opencensus.io v0.19.2
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
