module github.com/future-architect/futureot/exporters/opencensus-go-exporter-p8s-pushgateway

go 1.12

require (
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v1.0.0
	go.opencensus.io v0.22.0
)

replace github.com/future-architect/futureot/exporters/opencensus-go-exporter-p8s-pushgateway => ../../exporters/opencensus-go-exporter-p8s-pushgateway
