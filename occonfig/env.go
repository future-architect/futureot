package occonfig

import (
	"strings"
)

func envArrayToMap(envs []string) map[string]string {
	result := make(map[string]string)
	for _, env := range envs {
		entry := strings.SplitN(env, "=", 2)
		result[entry[0]] = entry[1]
	}
	return result
}

func initByEnvMap(envs []string) (*Config, error) {
	result := &Config{
		TraceSampler: -1.0,
	}
	envMaps := envArrayToMap(envs)
	if serviceName, ok := envMaps["OC_SERVICE_NAME"]; ok {
		result.ServiceName = serviceName
	}
	if serviceUrl, ok := envMaps["OC_SERVICE_URL"]; ok {
		result.ServiceUrl = serviceUrl
	}
	if zpage, ok := envMaps["OC_ZPAGE"]; ok {
		result.ZPage = zpage
	}
	if configJson, ok := envMaps["OC_CONFIG_JSON"]; ok {
		result.ConfigFile = configJson
	}
	if honeycombKey, ok := envMaps["OC_HONEYCOMB_WRITE_KEY"]; ok {
		result.HoneycombKey = honeycombKey
	}
	if tracer, ok := envMaps["OC_TRACE_EXPORTER"]; ok {
		result.TraceExporter = tracer
	}
	s, err := SelectSampler(envMaps["OC_TRACE_SAMPLER"])
	if err != nil {
		return nil, err
	} else {
		result.TraceSampler = s
	}
	if tracer, ok := envMaps["OC_STATS_EXPORTER"]; ok {
		result.StatsExporter = tracer
	}
	return result, nil
}
