export class Config {
    serviceName: string = ""
    serviceUrl: string = ""
    zPage: string = ""
    configFile: string = ""
    traceExporter: string = ""
    traceSampler: number = -1
    statsExporter: string = ""
}

export enum Mode {
    Trace = 1,
    Stats = 2
}

type commandParser = () => Config;
let comanndParser: commandParser | null = null;

export function registerCommandLineParser(parser: commandParser) {
    comanndParser = parser;
}

export function getCommandLineParser() {
    return comanndParser;
}

function selectString(a: string, b: string): string {
    if (b === "") {
        return a
    }
    return b
}

function selectNumber(a: number, b: number): number {
    if (b < 0) {
        return a
    }
    return b
}

export function mergeConfigs(low: Config, high: Config): Config {
    const config = new Config()
    config.serviceName = selectString(low.serviceName, high.serviceName)
    config.serviceUrl = selectString(low.serviceUrl, high.serviceUrl)
    config.zPage = selectString(low.zPage, high.zPage)
    config.configFile = selectString(low.configFile, high.configFile)
    config.traceExporter = selectString(low.traceExporter, high.traceExporter)
    config.traceSampler = selectNumber(low.traceSampler, high.traceSampler)
    config.statsExporter = selectString(low.statsExporter, high.statsExporter)
    return config
}
