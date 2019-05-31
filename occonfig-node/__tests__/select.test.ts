import { selectTraceExporter, selectStatsExporter, ExporterType } from "../src/select";

describe("selectTraceExporter", () => {
    type TestCase = {
        source: string,
        type: ExporterType,
        host: string
    }
    const testcases: TestCase[] = [
        { source: "stackdriver://my-project-id", type: ExporterType.STACKDRIVER, host: "my-project-id" },
        { source: "sd://my-project-id", type: ExporterType.STACKDRIVER, host: "my-project-id" },
        { source: "jaeger://localhost:14268", type: ExporterType.JAEGER, host: "http://localhost:14268/api/traces" },
        { source: "jaeger://localhost", type: ExporterType.JAEGER, host: "http://localhost:14268/api/traces" },
        { source: "jaeger", type: ExporterType.JAEGER, host: "http://localhost:14268/api/traces" },
        { source: "zipkin://localhost:9411/api/v2/spans", type: ExporterType.ZIPKIN, host: "http://localhost:9411/api/v2/spans" },
        { source: "zipkin://localhost:9411", type: ExporterType.ZIPKIN, host: "http://localhost:9411/api/v2/spans" },
        { source: "zipkin://localhost", type: ExporterType.ZIPKIN, host: "http://localhost:9411/api/v2/spans" },
        { source: "zipkin://localhost/", type: ExporterType.ZIPKIN, host: "http://localhost:9411/" },
        { source: "zipkin", type: ExporterType.ZIPKIN, host: "http://localhost:9411/api/v2/spans" },
    ];
    for (const testcase of testcases) {
        it(testcase.source, () => {
            const exporter = selectTraceExporter(testcase.source);
            expect(exporter).toBeTruthy();
            if (exporter !== null) {
                expect(exporter.host).toEqual(testcase.host);
                expect(exporter.type).toEqual(testcase.type);
            }
        })
    }
});

describe("selectStatsExporter", () => {
    type TestCase = {
        source: string,
        type: ExporterType,
        host: string
    }
    const testcases: TestCase[] = [
        { source: "stackdriver://my-project-id", type: ExporterType.STACKDRIVER, host: "my-project-id" },
        { source: "sd://my-project-id", type: ExporterType.STACKDRIVER, host: "my-project-id" },
        { source: "prometheus://:8888", type: ExporterType.PROMETHEUS, host: "http://localhost:8888" },
        { source: "p8s://:8888", type: ExporterType.PROMETHEUS, host: "http://localhost:8888" },
    ];
    for (const testcase of testcases) {
        it(testcase.source, () => {
            const exporter = selectStatsExporter(testcase.source);
            expect(exporter).toBeTruthy();
            if (exporter !== null) {
                expect(exporter.host).toEqual(testcase.host);
                expect(exporter.type).toEqual(testcase.type);
            }
        })
    }
});
