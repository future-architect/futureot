package occonfig

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type ExporterType int

const (
	STACKDRIVER ExporterType = iota
	XRAY
	DATADOG
	JAEGER
	ZIPKIN
	HONEYCOMB
	PROMETHEUS
	GRAPHITE
	ZAP
)

type Exporter struct {
	Type ExporterType
	Host string
}

func SelectTraceExporter(host string) (*Exporter, error) {
	u, _ := url.Parse(host)
	switch u.Scheme {
	case "sd":
		fallthrough
	case "stackdriver":
		return &Exporter{
			Type: STACKDRIVER,
			Host: u.Host,
		}, nil
	case "xray":
		return &Exporter{
			Type: XRAY,
		}, nil
	case "dd":
		fallthrough
	case "datadog":
		{
			host := u.Hostname()
			port := u.Port()
			if host == "" {
				host = "localhost"
			}
			if port == "" {
				port = "8125"
			}
			return &Exporter{
				Type: DATADOG,
				Host: fmt.Sprintf("%s:%s", host, port),
			}, nil

		}
	case "jaeger":
		{
			host := u.Hostname()
			port := u.Port()
			path := u.Path
			if host == "" {
				host = "localhost"
			}
			if port == "" {
				port = "14268"
			}
			if path == "" {
				path = "/api/traces"
			}
			return &Exporter{
				Type: JAEGER,
				Host: fmt.Sprintf("http://%s:%s%s", host, port, path),
			}, nil
		}
	case "jeager":
		return nil, errors.New("Misspelling! jeager -> jaeger")
	case "zipkin":
		{
			host := u.Hostname()
			port := u.Port()
			path := u.Path
			if host == "" {
				host = "localhost"
			}
			if port == "" {
				port = "9411"
			}
			if path == "" {
				path = "/api/v2/spans"
			}
			return &Exporter{
				Type: ZIPKIN,
				Host: fmt.Sprintf("http://%s:%s%s", host, port, path),
			}, nil
		}
	case "": // no scheme
		switch u.Path {
		case "xray":
			return &Exporter{
				Type: XRAY,
			}, nil
		case "dd":
			fallthrough
		case "datadog":
			return &Exporter{
				Type: DATADOG,
				Host: "localhost:8125",
			}, nil
		case "honeycomb":
			return &Exporter{
				Type: HONEYCOMB,
			}, nil
		case "jaeger":
			return &Exporter{
				Type: JAEGER,
				Host: "http://localhost:14268/api/traces",
			}, nil
		case "zipkin":
			return &Exporter{
				Type: ZIPKIN,
				Host: "http://localhost:9411/api/v2/spans",
			}, nil
		case "zap":
			return &Exporter{
				Type: ZAP,
			}, nil
		}
	}
	return nil, errors.New("No exporter config found")
}

func SelectStatsExporter(host string) (*Exporter, error) {
	u, _ := url.Parse(host)
	switch u.Scheme {
	case "sd":
		fallthrough
	case "stackdriver":
		return &Exporter{
			Type: STACKDRIVER,
			Host: u.Host,
		}, nil
	case "dd":
		fallthrough
	case "datadog":
		{
			host := u.Hostname()
			port := u.Port()
			if host == "" {
				host = "localhost"
			}
			if port == "" {
				port = "8125"
			}
			return &Exporter{
				Type: DATADOG,
				Host: fmt.Sprintf("%s:%s", host, port),
			}, nil

		}
	case "prometheus":
		fallthrough
	case "p8s":
		{
			port := u.Port()
			if port == "" {
				port = "8888"
			}
			return &Exporter{
				Type: PROMETHEUS,
				Host: fmt.Sprintf("http://:%s", port),
			}, nil

		}
	case "graphite":
		{
			host := u.Hostname()
			port := u.Port()
			if host == "" {
				host = "localhost"
			}
			if port == "" {
				port = "2003"
			}
			return &Exporter{
				Type: GRAPHITE,
				Host: fmt.Sprintf("%s:%s", host, port),
			}, nil
		}
	case "": // no scheme
		switch u.Path {
		case "dd":
			fallthrough
		case "datadog":
			return &Exporter{
				Type: DATADOG,
				Host: "localhost:8125",
			}, nil
		case "graphite":
			return &Exporter{
				Type: GRAPHITE,
				Host: "localhost:2003",
			}, nil
		}
	}
	return nil, errors.New("No exporter config found")
}

func SelectSampler(s string) (float64, error) {
	switch s {
	case "always":
		return 1.0, nil
	case "never":
		return 0.0, nil
	case "":
		return -1, nil
	default:
		return strconv.ParseFloat(s, 64)
	}
}
