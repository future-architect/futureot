import * as tracing from '@opencensus/nodejs';
import { Command } from "commander";
import { Mode, useCommander, init } from "../src";
import { basename } from "path";

async function sleep(duration: number) {
    return new Promise<number>((resolve) => {
        setTimeout(resolve, duration);
    })
}


async function main() {
    const program = new Command(basename(process.argv[1]));
    program.version("0.0.1");
    useCommander(program, Mode.Trace|Mode.Stats).parse(process.argv);
    init(Mode.Trace|Mode.Stats);

    tracing.tracer.startRootSpan({name: "trace-example"}, async rootSpan => {
        try {
            for (let i = 0; i < 10; i++) {
                await doWork(i);
            }
        } catch (e) {
            console.log(e);
        } finally {
            rootSpan.end();
        }
    });
}

async function doWork(i: number) {
    console.log("doing heavy task:", i);
    const span = tracing.tracer.startChildSpan(`doWork ${i}`);
    try {
        await sleep(400);
        span.addAttribute("task", i);
        //span.addAnnotation(`doing task ${i}`);
        await sleep(100);
    } catch (e) {
        console.log(e);
    } finally {
        span.end();
    }
}

main();
