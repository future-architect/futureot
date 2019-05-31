import { Command } from "commander";
import { Config, Mode, registerCommandLineParser } from "./config";
import { selectSampler } from "./select";

export function useCommander(program: Command, mode: Mode): Command {
    program
        .option("--oc-service-name <name>", "Service name that appears in OpenCensus resulting page")
        .option("--oc-service-url <url>", "Service URL")
        .option("--oc-config-json <path>", "Config JSON file path")
        .option("--oc-zpage <url>", "ZPage in-process debug console url (e.g. http://:8888/debug");
    if (mode & Mode.Trace) {
        program
            .option("--oc-trace-exporter <url>", "Trace exporter (e.g. stackdriver://demo-project-id, jaeger://localhost:6831")
            .option("--oc-trace-sampler <prob>", "Trace sampling rate ('always'(default), 'never', '0-1'")
    }
    if (mode & Mode.Stats) {
        program
            .option("--oc-stats-exporter <url>", "Stats exporter (e.g. stackdriver://demo-project-id, prometheus://localhost:8888")
    }
    registerCommandLineParser((): Config => {
        return parseCommander(program);
    });
    return program;
}

export function parseCommander(program: Command): Config {
    const config = new Config();
    config.serviceName = program.ocServiceName;
    config.serviceUrl = program.ocServiceUrl;
    config.configFile = program.ocConfigJson;
    config.zPage = program.ocZpage;
    config.traceExporter = program.ocTraceExporter;
    config.traceSampler = selectSampler(program.ocTraceSampler);
    config.statsExporter = program.ocStatsExporter;
    return config;
}
