package occonfig

import (
	"flag"
)

type FlagResult struct {
	ServiceName   string
	ServiceUrl    string
	HoneycombKey  string
	ConfigFile    string
	TraceExporter string
	TraceSampler  string
	StatsExporter string
	ZPage         string
}

var defaultFlagResult *FlagResult

func UseFlag(mode Mode, flagset ...*flag.FlagSet) {
	if len(flagset) == 0 {
		defaultFlagResult = InitFlagSet(flag.CommandLine, mode)
	} else {
		defaultFlagResult = InitFlagSet(flagset[0], mode)
	}
}

func InitFlagSet(flagset *flag.FlagSet, mode Mode) *FlagResult {
	result := &FlagResult{}
	flagset.StringVar(
		&result.ServiceName, "oc-service-name", "",
		"Service name that appears in OpenCensus resulting page")
	flagset.StringVar(
		&result.ServiceUrl, "oc-service-url", "",
		"Service URL")
	flagset.StringVar(
		&result.ConfigFile, "oc-config-json", "",
		"Config JSON file path")
	flagset.StringVar(
		&result.ZPage, "oc-zpage", "",
		"ZPage in-process debug console url (e.g. http://:8888/debug")

	if mode&Trace == Trace {
		flagset.StringVar(
			&result.TraceExporter, "oc-trace-exporter", "",
			"OpenCensus trace setting (e.g. stackdriver://demo-project-id, jaeger://localhost:6831")
		flagset.StringVar(
			&result.TraceSampler, "oc-trace-sampler", "",
			"Trace sampling rate ('always', 'never', '0-1'")
		flagset.StringVar(
			&result.HoneycombKey, "oc-honeycomb-write-key", "",
			"Honeycomb.io write key or file path(file://) (it is needed when trace exporter is honeycomb)")
	}

	if mode&Stats == Stats {
		flagset.StringVar(
			&result.StatsExporter, "oc-stats-exporter", "",
			"OpenCensus stats setting (e.g. stackdriver://demo-project-id, prometheus://localhost:8888")
	}

	getConfigFromCommandLine = func() (*Config, error) {
		return parseFlagResult(result)
	}

	return result
}

func parseFlagResult(flagResult *FlagResult) (config *Config, err error) {
	config = &Config{
		ServiceName:   flagResult.ServiceName,
		ServiceUrl:    flagResult.ServiceUrl,
		ConfigFile:    flagResult.ConfigFile,
		ZPage:         flagResult.ZPage,
		TraceExporter: flagResult.TraceExporter,
		HoneycombKey:  flagResult.HoneycombKey,
		StatsExporter: flagResult.StatsExporter,
	}
	s, e := SelectSampler(flagResult.TraceSampler)
	if e != nil {
		config = nil
		err = e
	} else {
		config.TraceSampler = s
	}
	return
}
