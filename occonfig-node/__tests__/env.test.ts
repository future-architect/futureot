import { initByEnvMap } from "../src/env";

describe("init by env", () => {
    type TestCase = {
        name: string
        envs: { [key: string]: string }
        serviceName?: string
        serviceUrl?: string
        zPage?: string
        configFile?: string
        traceExporter?: string
        traceSampler?: number
        statsExporter?: string
    };

    const testCases: TestCase[] = [
        {
            name: "service-name test",
            envs: {OC_SERVICE_NAME: "my-service", HOME: "test"},
            serviceName: "my-service",
        },
        {
            name: "service-url test",
            envs: {OC_SERVICE_URL: "http://localhost:8080", HOME: "test"},
            serviceUrl: "http://localhost:8080",
        },
        {
            name: "zpage test",
            envs: {OC_ZPAGE: "http://:8888/debug", HOME: "test"},
            zPage: "http://:8888/debug",
        },
        {
            name: "config-json test",
            envs: {OC_CONFIG_JSON: "config.json", HOME: "test"},
            configFile: "config.json",
        },
        {
            name: "trace-exporter test",
            envs: {OC_TRACE_EXPORTER: "jaeger://localhost:6831", HOME: "test"},
            traceExporter: "jaeger://localhost:6831",
        },
        {
            name: "trace-sampler test (1)",
            envs: {OC_TRACE_SAMPLER: "never", HOME: "test"},
            traceSampler: 0.0,
        },
        {
            name: "trace-sampler test (2)",
            envs: {OC_TRACE_SAMPLER: "always", HOME: "test"},
            traceSampler: 1.0,
        },
        {
            name: "trace-sampler test (3)",
            envs: {OC_TRACE_SAMPLER: "0.25", HOME: "test"},
            traceSampler: 0.25,
        },
        {
            name: "stats-exporter test",
            envs: {OC_STATS_EXPORTER: "prometheus://:8888", HOME: "test"},
            statsExporter: "prometheus://:8888",
        },
    ];

    testCases.forEach((testcase) => {
        it(testcase.name, () => {
            const result = initByEnvMap(testcase.envs)
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
