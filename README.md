# futureot

[OpenCensus](https://opencensus.io/) helper functions.

In near future, it will support [OpenTelemetry](https://opentelemetry.io/) too.

```sh
$ go get github.com/future-architect/futureot/...
```

## Config Helper

Package ``occonfig`` helps initializing OpenCensus.
The application that imports this package can initialize OpenCensus exporter setting
via environment variables, command line arguments.

## Instruments

* Under constructing

## Exporters

* [opencensus-go-exporter-zap](https://github.com/future-architect/futureot/tree/master/exporters/opencensus-go-exporter-zap): Console exporter via [zap](https://godoc.org/go.uber.org/zap).