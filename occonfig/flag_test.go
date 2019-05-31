package occonfig

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFlagSet(t *testing.T) {
	testcases := []struct {
		Name          string
		Params        []string
		ServiceName   string
		ServiceUrl    string
		ConfigFile    string
		HoneycombKey  string
		TraceExporter string
		TraceSampler  float64
		StatsExporter string
		ZPage         string
	}{
		{
			Name:         "service-name test",
			Params:       []string{"-oc-service-name", "my-service", "etc", "etc"},
			ServiceName:  "my-service",
			TraceSampler: -1,
		},
		{
			Name:         "service-url test (1)",
			Params:       []string{"-oc-service-url", "test"},
			ServiceUrl:   "test",
			TraceSampler: -1,
		},
		{
			Name:         "service-url test (2)",
			Params:       []string{"-oc-service-url", "http://localhost:8080"},
			ServiceUrl:   "http://localhost:8080",
			TraceSampler: -1,
		},
		{
			Name:         "config-json test",
			Params:       []string{"-oc-config-json", "config.json", "etc", "etc"},
			ConfigFile:   "config.json",
			TraceSampler: -1,
		},
		{
			Name:         "zpage test",
			Params:       []string{"-oc-zpage", "http://:8888/debug", "etc", "etc"},
			ZPage:        "http://:8888/debug",
			TraceSampler: -1,
		},
		{
			Name:         "honeycomb-write-key test",
			Params:       []string{"-oc-honeycomb-write-key", "honeycomb.key", "etc", "etc"},
			HoneycombKey: "honeycomb.key",
			TraceSampler: -1,
		},
		{
			Name:          "trace-exporter test",
			Params:        []string{"-oc-trace-exporter", "jaeger://localhost:6831", "etc", "etc"},
			TraceExporter: "jaeger://localhost:6831",
			TraceSampler:  -1,
		},
		{
			Name:         "trace-sampler test (1)",
			Params:       []string{"--oc-trace-sampler", "never"},
			TraceSampler: 0.0,
		},
		{
			Name:         "trace-sampler test (2)",
			Params:       []string{"--oc-trace-sampler", "always"},
			TraceSampler: 1.0,
		},
		{
			Name:         "trace-sampler test (3)",
			Params:       []string{"--oc-trace-sampler", "0.25"},
			TraceSampler: 0.25,
		},
		{
			Name:          "stats-exporter test",
			Params:        []string{"-oc-stats-exporter", "prometheus://localhost:8888", "etc", "etc"},
			StatsExporter: "prometheus://localhost:8888",
			TraceSampler:  -1,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			flagSet := flag.NewFlagSet("default", flag.ContinueOnError)
			flagSetResult := InitFlagSet(flagSet, Trace|Stats)
			err := flagSet.Parse(testcase.Params)
			result, err := parseFlagResult(flagSetResult)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, testcase.ServiceName, result.ServiceName)
			assert.Equal(t, testcase.ServiceUrl, result.ServiceUrl)
			assert.Equal(t, testcase.ConfigFile, result.ConfigFile)
			assert.Equal(t, testcase.ZPage, result.ZPage)
			assert.Equal(t, testcase.HoneycombKey, result.HoneycombKey)
			assert.Equal(t, testcase.TraceExporter, result.TraceExporter)
			assert.InDelta(t, testcase.TraceSampler, result.TraceSampler, 0.01)
			assert.Equal(t, testcase.StatsExporter, result.StatsExporter)
		})
	}
}
