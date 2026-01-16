# TypeGo

TypeGo is a TypeScript runtime built on Go. It lets you write TypeScript code that can directly import and use Go packages, leverage NPM modules, and compile everything into a standalone native binary.

## Why TypeGo?

TypeScript developers often need functionality that's easier to implement in Go—things like efficient file handling, network protocols, or leveraging existing Go libraries. TypeGo bridges that gap by embedding a JavaScript engine (Goja) into Go, then providing a seamless import system that connects the two worlds.

## Installation

Requires Go 1.21 or later.

```bash
go install github.com/repyh3/typego/cmd/typego@latest
```

Verify the installation:

```bash
typego --version
```

## Getting Started

Create a new project with the recommended structure and configuration:

```bash
typego init my-project
cd my-project
```

This generates a `tsconfig.json`, type definitions, and a sample `src/index.ts` file. Run your code with:

```bash
typego run src/index.ts
```

## Commands

| Command | Description |
|---------|-------------|
| `typego run <file>` | Execute a TypeScript file immediately |
| `typego build <file> -o <output>` | Compile to a standalone executable |
| `typego types` | Regenerate type definitions for Go imports |
| `typego init <name>` | Scaffold a new project |

## Importing Go Packages

TypeGo uses a special `go:` prefix to import Go packages. The runtime dynamically links these at compile time.

```typescript
import { Println, Printf } from "go:fmt";
import { Red, Green } from "go:github.com/fatih/color";

Println("Hello from Go's fmt package");
Red("Colored output from the fatih/color library");
```

When you import a third-party Go package, TypeGo fetches it using `go get` and generates the necessary bindings automatically.

## NPM Package Support

Standard NPM packages work out of the box. TypeGo bundles them using esbuild before execution.

```typescript
import _ from "lodash";
import chalk from "chalk";

const message = _.kebabCase("Hello TypeGo");
console.log(chalk.blue(message));
```

Run `npm install` in your project directory to add packages as usual.

## Workers

TypeGo supports background workers for concurrent execution. Each worker runs in its own isolated context.

```typescript
const worker = new Worker("worker.ts");

worker.onmessage = (event) => {
    console.log("Response from worker:", event.data);
};

worker.postMessage({ task: "compute", value: 42 });
```

Inside `worker.ts`:

```typescript
onmessage = (event) => {
    const result = event.data.value * 2;
    postMessage({ result });
};
```

## Shared Memory

For scenarios requiring shared state between the main thread and workers, TypeGo provides `makeShared`:

```typescript
const buffer = makeShared("counter", 4);
const view = new Int32Array(buffer);

// Atomic operations work as expected
Atomics.store(view, 0, 0);
Atomics.add(view, 0, 1);
```

## Node.js Compatibility

TypeGo includes polyfills for common Node.js globals so that many NPM packages work without modification:

- **process**: `process.env`, `process.platform`, `process.cwd()`, `process.argv`
- **Buffer**: `Buffer.from()`, `Buffer.alloc()`
- **Timers**: `setTimeout`, `setInterval`, `clearInterval`

## Building Standalone Binaries

Compile your TypeScript project into a single executable that runs anywhere Go supports:

```bash
typego build src/index.ts -o myapp.exe
```

The output binary includes the bundled JavaScript, the Go runtime, and all imported Go packages. No external dependencies required.

## License

MIT
