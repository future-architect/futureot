import { basename } from "path";

import { globalStats } from '@opencensus/core';
import tracing from '@opencensus/nodejs';
import { StackdriverStatsExporter, StackdriverTraceExporter } from '@opencensus/exporter-stackdriver';
import { JaegerTraceExporter } from '@opencensus/exporter-jaeger';
import { PrometheusStatsExporter } from '@opencensus/exporter-prometheus';
import { ZipkinTraceExporter } from '@opencensus/exporter-zipkin';

export { Mode } from "./config";
export { useCommander } from "./commander";

import { Mode, Config, getCommandLineParser, mergeConfigs } from "./config";
import { initByEnvMap } from "./env";
import { readFiles } from "./json";
import { selectTraceExporter, selectStatsExporter, ExporterType } from "./select";

function getConfig(): Config {
    let config = initByEnvMap(process.env as { [key: string]: string })
    config = readFiles(config, process.cwd());
    const parser = getCommandLineParser();
    if (parser) {
        let commandConfig = parser();
        commandConfig = readFiles(commandConfig, process.cwd());
        config = mergeConfigs(config, commandConfig);
    }
    if (config.traceSampler < 0) {
        config.traceSampler = 1.0
    }
    if (!config.serviceName) {
        config.serviceName = basename(process.argv[1]);
    }
    return config;
}

export function init(mode: Mode) {
    const config = getConfig();

    if (mode & Mode.Trace) {
        const exporter = selectTraceExporter(config.traceExporter);
        if (exporter) {
            switch (exporter.type) {
                case ExporterType.STACKDRIVER:
                    tracing.registerExporter(new StackdriverTraceExporter({
                        projectId: exporter.host
                    })).start();
                    break;
                case ExporterType.JAEGER: {
                    const url = new URL(exporter.host);
                    const jaegerExporter = new JaegerTraceExporter({
                        serviceName: config.serviceName,
                        host: url.hostname,
                        port: parseInt(url.port)
                    });
                    tracing.start({
                        exporter: jaegerExporter,
                    })
                    break;
                }
                case ExporterType.ZIPKIN:
                    tracing.registerExporter(new ZipkinTraceExporter({
                        serviceName: config.serviceName,
                        url: exporter.host,
                    })).start()
                    break;
            }
        } else {
            console.warn("[OpenCensus] No trace exporter found");
        }
    }

    if (mode & Mode.Stats) {
        const exporter = selectStatsExporter(config.statsExporter);
        if (exporter) {
            switch (exporter.type) {
                case ExporterType.STACKDRIVER:
                    globalStats.registerExporter(new StackdriverStatsExporter({
                        projectId: exporter.host
                    }));
                    break;
                case ExporterType.PROMETHEUS: {
                    const url = new URL(exporter.host);
                    globalStats.registerExporter(new PrometheusStatsExporter({
                        port: parseInt(url.port),
                        startServer: true
                    }));
                    break;
                }
            }
        } else {
            console.warn("[OpenCensus] No stats exporter found");
        }
    }
}

