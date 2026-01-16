import { Spawn, Sleep, Chan } from "go/sync";
import { Fetch } from "go/net/http";
import { Println } from "go/fmt";

async function main() {
    try {
        Println("🚀 Starting Advanced Demo...");
        Println(`Runtime: ${native.GetRuntimeInfo()} (Launched: ${native.StartTime})`);

        const alerts = new Chan<string>();

        Spawn(async () => {
            Println("[Sensor] Background worker started.");
            let cycle = 0;
            while (cycle < 3) {
                await Sleep(1500);
                cycle++;

                const logEntry = `Cycle ${cycle} anomaly detected.`;
                for (let i = 0; i < logEntry.length; i++) {
                    cliBuffer[i] = logEntry.charCodeAt(i);
                }
                cliBuffer[logEntry.length] = 0;

                alerts.send(`Alert #${cycle}`);
                Println(`[Sensor] Signal ${cycle} sent.`);
            }
        });

        Println("[Main] Waiting for alerts...");
        for (let i = 0; i < 3; i++) {
            const signal = await alerts.recv();
            Println(`[Main] Received: ${signal}`);

            let logData = "";
            for (let j = 0; j < 100; j++) {
                if (cliBuffer[j] === 0) break;
                logData += String.fromCharCode(cliBuffer[j]);
            }
            Println(`[Main] Shared Buffer Data: ${logData}`);

            Println("[Main] Calling Fetch API...");
            const resp = await Fetch("https://httpbin.org/get?typego=advanced");
            Println(`[Main] HTTP Task Completed. Status: ${resp.Status}`);
        }

        Println("✅ Advanced Demo Mission Accomplished.");
    } catch (e) {
        Println(`[FATAL] ${e}`);
    }
}

main();
