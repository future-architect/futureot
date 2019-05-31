module trace-sample

go 1.12

replace (
	github.com/future-architect/futureot/exporters/opencensus-go-exporter-zap => ../../exporters/opencensus-go-exporter-zap
	github.com/future-architect/futureot/occonfig => ../../occonfig
)

require (
	github.com/future-architect/futureot/occonfig v0.0.0-00010101000000-000000000000
	go.opencensus.io v0.22.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
