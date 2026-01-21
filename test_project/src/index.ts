import { Default } from "go:github.com/gin-gonic/gin";
import { Println } from "go:fmt";

async function main() {
    Println("🚀 Starting Gin Web Server via TypeGo JIT...");

    const r = Default();

    // In TypeGo, Go methods are mapped to JS properties/methods.
    // For this demo, we'll just initialize and run it.
    // Note: complex callbacks might need more setup, 
    // but running the engine is the first step.

    Println("📡 Listening on :8080");
    r.Run(":8080");
}

main();
