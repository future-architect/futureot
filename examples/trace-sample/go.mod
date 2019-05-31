module trace-sample

go 1.12

replace github.com/future-architect/futureot/occonfig => ../../occonfig

require (
	contrib.go.opencensus.io/exporter/jaeger v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/zipkin v0.1.1 // indirect
	github.com/future-architect/futureot/occonfig v0.0.0-00010101000000-000000000000
	go.opencensus.io v0.22.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
