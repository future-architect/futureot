package occonfig

import (
	"net/url"

	"gopkg.in/alecthomas/kingpin.v2"
)

type KingpinResult struct {
	ServiceName   string
	ServiceUrl    *url.URL
	HoneycombKey  string
	ConfigFile    string
	ZPage         *url.URL
	TraceExporter *url.URL
	TraceSampler  string
	StatsExporter *url.URL
}

func UseKingpin(mode Mode, application ...*kingpin.Application) {
	if len(application) == 0 {
		InitApplication(kingpin.CommandLine, mode)
	} else {
		InitApplication(application[0], mode)
	}
}

func InitApplication(application *kingpin.Application, mode Mode) *KingpinResult {
	result := &KingpinResult{}
	application.Flag("oc-service-name", "Service name that appears in OpenCensus resulting page").
		StringVar(&result.ServiceName)
	application.Flag("oc-service-url", "Service URL").
		URLVar(&result.ServiceUrl)
	application.Flag("oc-config-json", "Config JSON file path").
		ExistingFileVar(&result.ConfigFile)
	application.Flag("oc-zpage", "ZPage in-process debug console url (e.g. http://:8888/debug").
		URLVar(&result.ZPage)

	if mode&Trace == Trace {
		application.Flag("oc-trace-exporter", "Trace exporter (e.g. stackdriver://demo-project-id, jaeger://localhost:6831").
			URLVar(&result.TraceExporter)
		application.Flag("oc-trace-sampler", "Trace sampling rate ('always'(default), 'never', '0-1'").
			StringVar(&result.TraceSampler)
		application.Flag("oc-honeycomb-write-key", "Honeycomb.io write key or file path(file://) (it is needed when trace exporter is honeycomb)").
			ExistingFileVar(&result.HoneycombKey)
	}

	if mode&Stats == Stats {
		application.Flag("oc-stats-exporter", "Stats exporter (e.g. stackdriver://demo-project-id, prometheus://localhost:8888").
			URLVar(&result.StatsExporter)
	}

	getConfigFromCommandLine = func() (*Config, error) {
		return parseKingpinResult(result)
	}

	return result
}

func parseKingpinResult(kingpinResult *KingpinResult) (config *Config, err error) {
	config = &Config{
		ServiceName:  kingpinResult.ServiceName,
		HoneycombKey: kingpinResult.HoneycombKey,
		ConfigFile:   kingpinResult.ConfigFile,
	}
	if kingpinResult.ServiceUrl != nil {
		config.ServiceUrl = kingpinResult.ServiceUrl.String()
	}
	if kingpinResult.ZPage != nil {
		config.ZPage = kingpinResult.ZPage.String()
	}
	if kingpinResult.TraceExporter != nil {
		config.TraceExporter = kingpinResult.TraceExporter.String()
	}
	s, e := SelectSampler(kingpinResult.TraceSampler)
	if e != nil {
		config = nil
		err = e
	} else {
		config.TraceSampler = s
	}
	if kingpinResult.StatsExporter != nil {
		config.StatsExporter = kingpinResult.StatsExporter.String()
	}
	return
}
