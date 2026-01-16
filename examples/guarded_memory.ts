import { makeShared } from "go/memory";
import { Spawn, Sleep } from "go/sync";
import { Println } from "go/fmt";

async function main() {
    Println("🚀 Starting Guarded Memory Demo...");

    // 1. Create a 1KB shared buffer
    const shared = makeShared("telemetry", 1024);
    const { buffer, mutex } = shared;

    // 2. Writer Task: Periodic updates with Write Lock
    Spawn(async () => {
        let count = 0;
        while (count < 5) {
            await Sleep(1000);
            count++;

            Println(`[Writer] Requesting Write Lock...`);
            await mutex.lock(); // Async Promise
            try {
                const msg = `UPDATE_${count}_${new Date().toLocaleTimeString()}`;
                for (let i = 0; i < msg.length; i++) buffer[i] = msg.charCodeAt(i);
                buffer[msg.length] = 0;
                Println(`[Writer] Data committed: ${msg}`);
            } finally {
                mutex.unlock();
                Println(`[Writer] Released Write Lock.`);
            }
        }
        Println("[Writer] Job finished.");
    });

    // 3. Reader Tasks: Concurrent reading with Read Lock
    const startReader = (id: number) => {
        Spawn(async () => {
            for (let i = 0; i < 5; i++) {
                await Sleep(Math.random() * 2000);

                Println(`[Reader ${id}] Requesting Read Lock...`);
                await mutex.rlock(); // Concurrent-friendly lock
                try {
                    let data = "";
                    for (let j = 0; j < buffer.length; j++) {
                        if (buffer[j] === 0) break;
                        data += String.fromCharCode(buffer[j]);
                    }
                    Println(`[Reader ${id}] Read: ${data}`);
                } finally {
                    mutex.runlock();
                    Println(`[Reader ${id}] Released Read Lock.`);
                }
            }
            Println(`[Reader ${id}] Job finished.`);
        });
    };

    startReader(1);
    startReader(2);

    Println("[Main] Concurrency tasks orchestrated.");
}

main();
