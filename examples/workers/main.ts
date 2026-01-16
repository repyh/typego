import { Println } from "go/fmt";
import { Sleep, Spawn } from "go/sync";
import { makeShared } from "go/memory";

async function main() {
    Println("🚀 TypeGo Worker Cluster Demo");
    Println("Initializing shared memory...");

    // Create shared buffer
    const shared = makeShared("cluster_data", 1024);
    shared.buffer[0] = 0;

    const workers: Worker[] = [];
    const NUM_WORKERS = 4;

    Println(`Spawning ${NUM_WORKERS} workers...`);

    // Spawn Workers
    for (let i = 0; i < NUM_WORKERS; i++) {
        const w = new Worker("examples/workers/worker.ts");

        w.onmessage = (e) => {
            Println(`[Main] Received from Worker ${i}: ${JSON.stringify(e.data)}`);
            if (e.data.val === 20) {
                Println("🎉 Target reached! Terminating cluster.");
                workers.forEach(worker => worker.terminate());
            }
        };

        workers.push(w);
    }

    await Sleep(1000);

    // Send commands
    Spawn(async () => {
        for (let i = 0; i < 20; i++) {
            const workerParams = i % NUM_WORKERS;
            Println(`[Main] Sending increment task to Worker ${workerParams}`);
            workers[workerParams].postMessage({ cmd: "increment" });
            await Sleep(100);
        }
    });

    // Keep alive for demo
    await Sleep(5000);
    Println("Demo finished.");
}

main();
