import { Command } from "commander";
import { Mode } from "../src/config";
import { useCommander, parseCommander } from "../src/commander";

describe("init by commander", () => {
    type TestCase = {
        name: string
        argv: string[]
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
            argv: ["--oc-service-name=my-service"],
            serviceName: "my-service",
        },
        {
            name: "service-url test",
            argv: ["--oc-service-url=test"],
            serviceUrl: "test",
        },
        {
            name: "service-url test (2)",
            argv: ["--oc-service-url=http://localhost:8080"],
            serviceUrl: "http://localhost:8080",
        },
        {
            name: "config-json test",
            argv: ["--oc-config-json", "./testdata/config.json"],
            configFile: "./testdata/config.json",
        },
        {
            name: "zpage test",
            argv: ["--oc-zpage", "http://:8888/debug"],
            zPage: "http://:8888/debug",
        },
        {
            name: "trace-exporter test",
            argv: ["--oc-trace-exporter", "jaeger://localhost:6831"],
            traceExporter: "jaeger://localhost:6831",
        },
        {
            name: "trace-sampler test (1)",
            argv: ["--oc-trace-sampler", "never"],
            traceSampler: 0.0,
        },
        {
            name: "trace-sampler test (2)",
            argv: ["--oc-trace-sampler", "always"],
            traceSampler: 1.0,
        },
        {
            name: "trace-sampler test (3)",
            argv: ["--oc-trace-sampler", "0.25"],
            traceSampler: 0.25,
        },
        {
            name: "stats-exporter test",
            argv: ["--oc-stats-exporter", "prometheus://localhost:8888"],
            statsExporter: "prometheus://localhost:8888",
        },
    ];

    let mockExit: any;

    beforeEach(() => {
        mockExit = jest.spyOn(process, 'exit').mockImplementation((code?: number): never => {
            throw new Error(`error: ${code}`);
        });
    })

    afterEach(() => {
        mockExit.mockRestore();
    })

    testCases.forEach((testcase) => {
        it(testcase.name, () => {
            const program = new Command();
            useCommander(program
                .version('0.0.1'), Mode.Trace | Mode.Stats)
                .parse(["node", "index.js", ...testcase.argv]);
            const config = parseCommander(program);
            if (testcase.serviceName) {
                expect(config.serviceName).toEqual(testcase.serviceName);
            }
            if (testcase.serviceUrl) {
                expect(config.serviceUrl).toEqual(testcase.serviceUrl);
            }
            if (testcase.zPage) {
                expect(config.zPage).toEqual(testcase.zPage);
            }
            if (testcase.configFile) {
                expect(config.configFile).toEqual(testcase.configFile);
            }
            if (testcase.traceExporter) {
                expect(config.traceExporter).toEqual(testcase.traceExporter);
            }
            if (testcase.traceSampler) {
                expect(config.traceSampler).toBeCloseTo(testcase.traceSampler);
            }
            if (testcase.statsExporter) {
                expect(config.statsExporter).toEqual(testcase.statsExporter);
            }
        });
    });
});
