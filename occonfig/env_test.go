package occonfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	testcases := []struct {
		Name          string
		Envs          []string
		ServiceName   string
		ServiceUrl    string
		ZPage         string
		ConfigFile    string
		HoneycombKey  string
		TraceExporter string
		TraceSampler  float64
		StatsExporter string
	}{
		{
			Name:         "service-name test",
			Envs:         []string{"OC_SERVICE_NAME=my-service", "HOME=test"},
			ServiceName:  "my-service",
			TraceSampler: -1,
		},
		{
			Name:         "service-url test",
			Envs:         []string{"OC_SERVICE_URL=http://localhost:8080", "HOME=test"},
			ServiceUrl:   "http://localhost:8080",
			TraceSampler: -1,
		},
		{
			Name:         "zpage test",
			Envs:         []string{"OC_ZPAGE=http://:8888/debug", "HOME=test"},
			ZPage:        "http://:8888/debug",
			TraceSampler: -1,
		},
		{
			Name:         "config-json test",
			Envs:         []string{"OC_CONFIG_JSON=config.json", "HOME=test"},
			ConfigFile:   "config.json",
			TraceSampler: -1,
		},
		{
			Name:         "honeycomb-write-key test",
			Envs:         []string{"OC_HONEYCOMB_WRITE_KEY=honeycomb.key", "HOME=test"},
			HoneycombKey: "honeycomb.key",
			TraceSampler: -1,
		},
		{
			Name:          "trace-exporter test",
			Envs:          []string{"OC_TRACE_EXPORTER=jaeger://localhost:6831", "HOME=test"},
			TraceExporter: "jaeger://localhost:6831",
			TraceSampler:  -1,
		},
		{
			Name:         "trace-sampler test (1)",
			Envs:         []string{"OC_TRACE_SAMPLER=never", "HOME=test"},
			TraceSampler: 0.0,
		},
		{
			Name:         "trace-sampler test (2)",
			Envs:         []string{"OC_TRACE_SAMPLER=always", "HOME=test"},
			TraceSampler: 1.0,
		},
		{
			Name:         "trace-sampler test (3)",
			Envs:         []string{"OC_TRACE_SAMPLER=0.25", "HOME=test"},
			TraceSampler: 0.25,
		},
		{
			Name:          "stats-exporter test",
			Envs:          []string{"OC_STATS_EXPORTER=prometheus://:8888", "HOME=test"},
			StatsExporter: "prometheus://:8888",
			TraceSampler:  -1,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			result, err := initByEnvMap(testcase.Envs)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, testcase.ServiceName, result.ServiceName)
			assert.Equal(t, testcase.ServiceUrl, result.ServiceUrl)
			assert.Equal(t, testcase.ZPage, result.ZPage)
			assert.Equal(t, testcase.ConfigFile, result.ConfigFile)
			assert.Equal(t, testcase.HoneycombKey, result.HoneycombKey)
			assert.Equal(t, testcase.TraceExporter, result.TraceExporter)
			assert.Equal(t, testcase.StatsExporter, result.StatsExporter)
			assert.InDelta(t, testcase.TraceSampler, result.TraceSampler, 0.01)
		})
	}
}
