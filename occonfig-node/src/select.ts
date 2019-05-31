export enum ExporterType {
    STACKDRIVER,
    JAEGER,
    ZIPKIN,
    PROMETHEUS,
}

export type Exporter = { type: ExporterType, host: string };

export function selectTraceExporter(s: string): Exporter | null {
    switch (s) {
        case "jaeger":
            return {
                type: ExporterType.JAEGER,
                host: "http://localhost:6832",
            }
        case "zipkin":
            return {
                type: ExporterType.ZIPKIN,
                host: "http://localhost:9411/api/v2/spans",
            }
        case "":
            return null;
    }
    const url = new URL(s);
    switch (url.protocol) {
        case "sd:":
        case "stackdriver:":
            return {
                type: ExporterType.STACKDRIVER,
                host: url.hostname,
            }
        case "jaeger:":
            if (url.hostname === "") {
                url.hostname = "localhost"
            }
            if (url.port == "") {
                url.port = "6832"
            }
            return {
                type: ExporterType.JAEGER,
                host: `http://${url.hostname}:${url.port}${url.pathname}`,
            }
        case "jeager:":
            throw new Error("Misspelling! jeager -> jaeger")
        case "zipkin:":
            if (url.hostname === "") {
                url.hostname = "localhost"
            }
            if (url.port == "") {
                url.port = "9411"
            }
            if (url.pathname == "") {
                url.pathname = "/api/v2/spans"
            }
            return {
                type: ExporterType.ZIPKIN,
                host: `http://${url.hostname}:${url.port}${url.pathname}`,
            }
    }
    return null;
}

export function selectStatsExporter(s: string): Exporter | null {
    if (!s) {
        return null;
    }
    const match = /([a-zA-Z0-9]+):\/\/:([0-9]{1,4})/.exec(s);
    if (match) {
        s = `${match[1]}://localhsot${match[2]}`;
    }
    const url = new URL(s);
    switch (url.protocol) {
        case "sd:":
        case "stackdriver:":
            return {
                type: ExporterType.STACKDRIVER,
                host: url.hostname,
            }
        case "prometheus:":
        case "p8s:":
            if (url.port == "") {
                url.port = "8888"
            }
            return {
                type: ExporterType.PROMETHEUS,
                host: `http://localhost:${url.port}`,
            }
    }
    return null;
}

export function selectSampler(s: string): number {
    switch (s) {
        case "always":
            return 1.0
        case "never":
            return 0.0
        case "":
        case undefined:
            return -1
        default:
            return parseFloat(s)
    }
}
