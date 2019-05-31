import { Config } from "./config";
import { selectSampler } from "./select";

export function initByEnvMap(env: { [key: string]: string }): Config {
    const result = new Config();
    for (const [key, value] of Object.entries(env)) {
        switch (key) {
            case "OC_SERVICE_NAME":
                result.serviceName = value
                break;
            case "OC_SERVICE_URL":
                result.serviceUrl = value
                break;
            case "OC_ZPAGE":
                result.zPage = value
                break;
            case "OC_CONFIG_JSON":
                result.configFile = value
                break;
            case "OC_TRACE_EXPORTER":
                result.traceExporter = value
                break;
            case "OC_TRACE_SAMPLER":
                result.traceSampler = selectSampler(value)
                break;
            case "OC_STATS_EXPORTER":
                result.statsExporter = value
                break;
        }
    }
    return result
}
