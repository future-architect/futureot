module stats_sample

go 1.12

replace (
	github.com/future-architect/futureot/exporters/opencensus-go-exporter-p8s-pushgateway => ../../exporters/opencensus-go-exporter-p8s-pushgateway
	github.com/future-architect/futureot/exporters/opencensus-go-exporter-zap => ../../exporters/opencensus-go-exporter-zap
	github.com/future-architect/futureot/occonfig => ../../occonfig
)

require (
	github.com/future-architect/futureot/exporters/opencensus-go-exporter-p8s-pushgateway v0.0.0-00010101000000-000000000000
	go.opencensus.io v0.22.0
)
