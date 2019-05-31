package occonfig

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseJson(t *testing.T) {
	testcases := []struct {
		Name          string
		Source        string
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
			Name:         "serviceName test",
			Source:       `{"serviceName": "my-service"}`,
			ServiceName:  "my-service",
			TraceSampler: -1,
		},
		{
			Name:         "serviceUrl test (1)",
			Source:       `{"serviceUrl": "test"}`,
			ServiceUrl:   "test",
			TraceSampler: -1,
		},
		{
			Name:         "serviceUrl test (2)",
			Source:       `{"serviceUrl": "http://localhost:8080"}`,
			ServiceUrl:   "http://localhost:8080",
			TraceSampler: -1,
		},
		{
			Name:         "zpage test",
			Source:       `{"zpage": "http://:8080/debug"}`,
			ZPage:        "http://:8080/debug",
			TraceSampler: -1,
		},
		{
			Name:         "configJson test",
			Source:       `{"extends": "./testdata/config.json"}`,
			ConfigFile:   "./testdata/config.json",
			TraceSampler: -1,
		},
		{
			Name:         "honeycombKey test",
			Source:       `{"trace": {"honeycombWriteKey": "./testdata/honeycomb.key"} }`,
			HoneycombKey: "./testdata/honeycomb.key",
			TraceSampler: -1,
		},
		{
			Name:          "trace-exporter test",
			Source:        `{"trace": {"exporter": "jaeger://localhost:6831"} }`,
			TraceExporter: "jaeger://localhost:6831",
			TraceSampler:  -1,
		},
		{
			Name:         "trace-sampler test (1)",
			Source:       `{"trace": {"sampler": "never"} }`,
			TraceSampler: 0.0,
		},
		{
			Name:         "trace-sampler test (2)",
			Source:       `{"trace": {"sampler": "always"} }`,
			TraceSampler: 1,
		},
		{
			Name:         "trace-sampler test (3)",
			Source:       `{"trace": {"sampler": 0.25} }`,
			TraceSampler: 0.25,
		},
		{
			Name:          "stats-exporter test",
			Source:        `{"stats": {"exporter": "p8s://localhost:8888"} }`,
			StatsExporter: "p8s://localhost:8888",
			TraceSampler:  -1,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			result, err := parseJSON([]byte(testcase.Source))
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
			assert.InDelta(t, testcase.TraceSampler, result.TraceSampler, 0.01)
			assert.Equal(t, testcase.StatsExporter, result.StatsExporter)
		})
	}
}

func TestReadJsonConfig1(t *testing.T) {
	config := &Config{
		ConfigFile: "./testdata/config.json",
	}
	wd, _ := os.Getwd()
	config, err := readFiles(config, wd)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, "my-service-name-at-file", config.ServiceName)
	assert.Equal(t, "dummy key at file", config.HoneycombKey)
	assert.Equal(t, "stackdriver://demo-project-id", config.TraceExporter)
}

func TestReadJsonConfig2(t *testing.T) {
	config := &Config{
		ServiceName:  "my-service-name-at-config",
		ServiceUrl:   "http://url-at-config",
		HoneycombKey: "dummy at config",
		ConfigFile:   "./testdata/config.json",
	}
	wd, _ := os.Getwd()
	config, err := readFiles(config, wd)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, "my-service-name-at-config", config.ServiceName)
	assert.Equal(t, "http://url-at-config", config.ServiceUrl)
	assert.Equal(t, "dummy at config", config.HoneycombKey)
}

func TestReadJsonWithExtendsConfig(t *testing.T) {
	config := &Config{
		ConfigFile: "./testdata/extends.json",
	}
	wd, _ := os.Getwd()
	config, err := readFiles(config, wd)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, "my-service-name-at-extends-json", config.ServiceName)
}
