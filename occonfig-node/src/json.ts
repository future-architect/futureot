import { resolve, dirname } from "path";
import { readFileSync } from "fs"
import { Config, mergeConfigs } from "./config";

export function readFiles(config: Config, currentFolder: string): Config {
    while (config.configFile) {
        const filePath = resolve(currentFolder, config.configFile)
        currentFolder = dirname(filePath)
        const jsonConfig = parseJSON(readFileSync(filePath, "utf8"))
        config.configFile = ""
        config = mergeConfigs(jsonConfig, config)
    }
    return config
}

export function parseJSON(content: string): Config {
    const config = new Config();
    const root = JSON.parse(content);
    config.serviceName = root["serviceName"] ? root["serviceName"] : "";
    config.serviceUrl = root["serviceUrl"] ? root["serviceUrl"] : "";
    config.zPage = root["zpage"] ? root["zpage"] : "";
    config.configFile = root["extends"] ? root["extends"] : "";

    if (root.trace) {
        const trace = root.trace;
        config.configFile = trace["configFile"] ? trace["configFile"] : "";
        config.traceExporter = trace["exporter"] ? trace["exporter"] : "";
        const rawSampler = trace["sampler"];
        switch (rawSampler) {
            case "always":
                config.traceSampler = 1.0;
                break;
            case "never":
                config.traceSampler = 0.0;
                break;
            case undefined:
            case "":
                break;
            default:
                config.traceSampler = parseFloat(rawSampler);
        }
    }
    if (root.stats) {
        config.statsExporter = root.stats["exporter"] ? root.stats["exporter"] : "";
    }
    return config
}
