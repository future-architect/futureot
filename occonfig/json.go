package occonfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func getString(tree map[string]interface{}, key string) string {
	if rawValue, ok := tree[key]; ok {
		if strValue, ok := rawValue.(string); ok {
			return strValue
		}
	}
	return ""
}

func parseJSON(content []byte) (config *Config, err error) {
	config = &Config{
		TraceSampler: -1,
	}
	root := make(map[string]interface{})
	err = json.Unmarshal(content, &root)
	if err != nil {
		config = nil
		return
	}
	config = &Config{
		ServiceName:  getString(root, "serviceName"),
		ServiceUrl:   getString(root, "serviceUrl"),
		ZPage:        getString(root, "zpage"),
		ConfigFile:   getString(root, "extends"),
		TraceSampler: -1,
	}
	if rawTrace, ok := root["trace"]; ok {
		if trace, ok := rawTrace.(map[string]interface{}); ok {
			config.HoneycombKey = getString(trace, "honeycombWriteKey")
			config.TraceExporter = getString(trace, "exporter")
			if rawSampler, ok := trace["sampler"]; ok {
				switch value := rawSampler.(type) {
				case string:
					switch value {
					case "always":
						config.TraceSampler = 1.0
					case "never":
						config.TraceSampler = 0.0
					default:
						config = nil
						err = errors.New("Invalid value trace.sampler. It should be 'always'|'never'|floating number(0-1).")
					}
				case float64:
					config.TraceSampler = value
				}
			}
		}
	}
	if rawTrace, ok := root["stats"]; ok {
		if trace, ok := rawTrace.(map[string]interface{}); ok {
			config.StatsExporter = getString(trace, "exporter")
		}
	}
	return
}

func selectString(a, b string) string {
	if b == "" {
		return a
	}
	return b
}

func selectNumber(a, b float64) float64 {
	if b < 0 {
		return a
	}
	return b
}

func mergeConfigs(low, high *Config) *Config {
	return &Config{
		ServiceName:   selectString(low.ServiceName, high.ServiceName),
		ServiceUrl:    selectString(low.ServiceUrl, high.ServiceUrl),
		ZPage:         selectString(low.ZPage, high.ZPage),
		ConfigFile:    selectString(low.ConfigFile, high.ConfigFile),
		TraceExporter: selectString(low.TraceExporter, high.TraceExporter),
		TraceSampler:  selectNumber(low.TraceSampler, high.TraceSampler),
		HoneycombKey:  selectString(low.HoneycombKey, high.HoneycombKey),
		StatsExporter: selectString(low.StatsExporter, high.StatsExporter),
	}
}

func readHoneycombKey(config *Config, currentFolder string) error {
	if strings.HasPrefix(config.HoneycombKey, "file://") {
		keyPath := strings.TrimPrefix(config.HoneycombKey, "file://")
		content, err := ioutil.ReadFile(filepath.Join(currentFolder, keyPath))
		if err != nil {
			return err
		}
		config.HoneycombKey = string(content)
	}
	return nil
}

func readFiles(config *Config, currentFolder string) (*Config, error) {
	err := readHoneycombKey(config, currentFolder)
	if err != nil {
		return nil, err
	}
	for config.ConfigFile != "" {
		filePath := filepath.Clean(filepath.Join(currentFolder, config.ConfigFile))
		currentFolder = filepath.Dir(filePath)
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		jsonConfig, err := parseJSON(file)
		if err != nil {
			return nil, err
		}
		err = readHoneycombKey(jsonConfig, currentFolder)
		if err != nil {
			return nil, err
		}
		config.ConfigFile = ""
		config = mergeConfigs(jsonConfig, config)
	}
	return config, nil
}
