import { Println } from "go/fmt";
import { Sleep } from "go/sync";

async function main() {
    Println("🚀 TypeGo Project Initialized!");
    await Sleep(500);
    Println("Happy coding!");
}

main();
