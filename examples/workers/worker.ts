import { Println } from "go/fmt";
import { Sleep, Mutex } from "go/sync";
import { makeShared } from "go/memory";

// Initialize the shared buffer and mutex
const shared = makeShared("cluster_data", 1024);
// Note: Workers will attach to this automatically via the same Name using makeShared internally

Println(`[Worker] Started.`);

self.onmessage = async (e) => {
    const msg = e.data;
    Println(`[Worker] Received: ${JSON.stringify(msg)}`);

    if (msg.cmd === "increment") {
        await shared.mutex.lock();

        let val = shared.buffer[0];
        val++;
        shared.buffer[0] = val;

        Println(`[Worker] Incremented to ${val}`);

        shared.mutex.unlock();

        self.postMessage({ type: "done", val: val });
    }
};
