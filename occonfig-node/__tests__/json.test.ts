import { parseJSON } from "../src/json";

describe("init by json", () => {
    interface TestCase {
        name: string;
        source: string;
        serviceName?: string;
        serviceUrl?: string;
        zPage?: string;
        configFile?: string;
        traceExporter?: string;
        traceSampler: number;
        statsExporter?: string;
    }

    const testCases: TestCase[] = [
        {
            name: "service-name test",
            source: `{"serviceName": "my-service"}`,
            serviceName: "my-service",
            traceSampler: -1
        },
        {
            name: "service-url test (1)",
            source: `{"serviceUrl": "test"}`,
            serviceUrl: "test",
            traceSampler: -1
        },
        {
            name: "service-url test (2)",
            source: `{"serviceUrl": "http://localhost:8080"}`,
            serviceUrl: "http://localhost:8080",
            traceSampler: -1
        },
        {
            name: "zpage test",
            source: `{"zpage": "http://:8080/debug"}`,
            zPage: "http://:8080/debug",
            traceSampler: -1
        },
        {
            name: "config-json test",
            source: `{"extends": "./testdata/config.json"}`,
            configFile: "./testdata/config.json",
            traceSampler: -1
        },
        {
            name: "trace-exporter test",
            source: `{"trace": {"exporter": "jaeger://localhost:6831"} }`,
            traceExporter: "jaeger://localhost:6831",
            traceSampler: -1
        },
        {
            name: "trace-sampler test (1)",
            source: `{"trace": {"sampler": "never"} }`,
            traceSampler: 0.0
        },
        {
            name: "trace-sampler test (2)",
            source: `{"trace": {"sampler": "always"} }`,
            traceSampler: 1.0
        },
        {
            name: "trace-sampler test (3)",
            source: `{"trace": {"sampler": 0.25} }`,
            traceSampler: 0.25
        },
        {
            name: "stats-exporter test",
            source: `{"stats": {"exporter": "p8s://localhost:8888"} }`,
            statsExporter: "p8s://localhost:8888",
            traceSampler: -1
        }
    ];
    testCases.forEach(testcase => {
        it(testcase.name, () => {
            const result = parseJSON(testcase.source);
            if (testcase.serviceName) {
                expect(result.serviceName).toEqual(testcase.serviceName);
            }
            if (testcase.serviceUrl) {
                expect(result.serviceUrl).toEqual(testcase.serviceUrl);
            }
            if (testcase.zPage) {
                expect(result.zPage).toEqual(testcase.zPage);
            }
            if (testcase.configFile) {
                expect(result.configFile).toEqual(testcase.configFile);
            }
            if (testcase.traceExporter) {
                expect(result.traceExporter).toEqual(testcase.traceExporter);
            }
            if (testcase.traceSampler) {
                expect(result.traceSampler).toBeCloseTo(testcase.traceSampler);
            }
            if (testcase.statsExporter) {
                expect(result.statsExporter).toEqual(testcase.statsExporter);
            }
        });
    });
});

/*
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
*/
