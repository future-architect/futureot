package occonfig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectTraceExporterOK(t *testing.T) {
	testcases := []struct {
		Source string
		Type   ExporterType
		Host   string
	}{
		{"stackdriver://my-project-id", STACKDRIVER, "my-project-id"},
		{"sd://my-project-id", STACKDRIVER, "my-project-id"},
		{"datadog://localhost:8125", DATADOG, "localhost:8125"},
		{"datadog://localhost", DATADOG, "localhost:8125"},
		{"datadog", DATADOG, "localhost:8125"},
		{"dd://localhost:8125", DATADOG, "localhost:8125"},
		{"dd://localhost", DATADOG, "localhost:8125"},
		{"dd", DATADOG, "localhost:8125"},
		{"xray", XRAY, ""},
		{"honeycomb", HONEYCOMB, ""},
		{"jaeger://localhost:14268", JAEGER, "http://localhost:14268/api/traces"},
		{"jaeger://localhost", JAEGER, "http://localhost:14268/api/traces"},
		{"jaeger", JAEGER, "http://localhost:14268/api/traces"},
		{"zipkin://localhost:9411/api/v2/spans", ZIPKIN, "http://localhost:9411/api/v2/spans"},
		{"zipkin://localhost:9411", ZIPKIN, "http://localhost:9411/api/v2/spans"},
		{"zipkin://localhost", ZIPKIN, "http://localhost:9411/api/v2/spans"},
		{"zipkin://localhost/", ZIPKIN, "http://localhost:9411/"},
		{"zipkin", ZIPKIN, "http://localhost:9411/api/v2/spans"},
		{"zap", ZAP, ""},
	}
	for _, testcase := range testcases {
		t.Run(testcase.Source, func(t *testing.T) {
			exporter, err := SelectTraceExporter(testcase.Source)
			assert.Nil(t, err)
			if exporter != nil {
				t.Log(testcase.Source, err)
				assert.Equal(t, testcase.Type, exporter.Type)
				assert.Equal(t, testcase.Host, exporter.Host)
			}
		})
	}
}

func TestSelectStatsExporterOK(t *testing.T) {
	testcases := []struct {
		Source string
		Type   ExporterType
		Host   string
	}{
		{"stackdriver://my-project-id", STACKDRIVER, "my-project-id"},
		{"sd://my-project-id", STACKDRIVER, "my-project-id"},
		{"datadog://localhost:8125", DATADOG, "localhost:8125"},
		{"datadog://localhost", DATADOG, "localhost:8125"},
		{"datadog", DATADOG, "localhost:8125"},
		{"dd://localhost:8125", DATADOG, "localhost:8125"},
		{"dd://localhost", DATADOG, "localhost:8125"},
		{"dd", DATADOG, "localhost:8125"},
		{"prometheus://:8888", PROMETHEUS, "http://:8888"},
		{"p8s://:8888", PROMETHEUS, "http://:8888"},
		{"graphite://:2003", GRAPHITE, "localhost:2003"},
		{"graphite", GRAPHITE, "localhost:2003"},
	}
	for _, testcase := range testcases {
		t.Run(testcase.Source, func(t *testing.T) {
			exporter, err := SelectStatsExporter(testcase.Source)
			assert.Nil(t, err)
			if exporter != nil {
				t.Log(testcase.Source, err)
				assert.Equal(t, testcase.Type, exporter.Type)
				assert.Equal(t, testcase.Host, exporter.Host)
			}
		})
	}
}
