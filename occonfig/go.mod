module github.com/future-architect/futureot/occonfig

go 1.12

replace github.com/future-architect/futureot/exporters/opencensus-go-exporter-zap => ../exporters/opencensus-go-exporter-zap

require (
	contrib.go.opencensus.io/exporter/aws v0.0.0-20181029163544-2befc13012d0
	contrib.go.opencensus.io/exporter/graphite v0.0.0-20190325161142-f4bcbbf058a5
	contrib.go.opencensus.io/exporter/jaeger v0.1.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	contrib.go.opencensus.io/exporter/stackdriver v0.9.2
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/DataDog/datadog-go v0.0.0-20190323183505-07c7c350327b // indirect
	github.com/Datadog/opencensus-go-exporter-datadog v0.0.0-20190314110122-1e6ba4554ec1
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/limitgroup v0.0.0-20150612190941-6abd8d71ec01 // indirect
	github.com/facebookgo/muster v0.0.0-20150708232844-fd3d7953fd52 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/future-architect/futureot/exporters/opencensus-go-exporter-zap v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/honeycombio/libhoney-go v1.9.5 // indirect
	github.com/honeycombio/opencensus-exporter v0.0.0-20190305003048-7d9e6ede23e2
	github.com/openzipkin/zipkin-go v0.1.6
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/tinylib/msgp v1.1.0 // indirect
	go.opencensus.io v0.22.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.11.0 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/alexcesaro/statsd.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
