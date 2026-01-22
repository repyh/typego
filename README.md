<div align="center">

# TypeGo

typego is an embedded TypeScript runtime for Go. It lets you script Go applications with TS without the overhead of Node.js or the boilerplate of manual FFI bindings.

[Getting Started](#getting-started) • [Features](#features) • [Examples](#examples) • [Optimization](OPTIMIZATION.md) • [Contributing](CONTRIBUTING.md) • [License](#license)

</div>

> [!NOTE]
> **Project Status**: TypeGo is under active development. However, please note that **maintenance is limited** as I am balancing this project with my university commitments. Issues and PRs are welcome but may see delayed responses.

Unlike typical runtimes that communicate over IPC or JSON-RPC, typego runs a JS engine (Sobek) directly inside your Go process. You can import Go packages as if they were native TS modules using the go: prefix, allowing for zero-copy data sharing and direct access to Go’s standard library.

## Index

- [Overview](#overview)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Usage](#usage)
- [Language Reference](#language-reference)
  - [Imports](#imports)
  - [Intrinsics](#intrinsics)
  - [Standard Library](#standard-library)
    - [go:fmt](#gofmt)
    - [go:net/http](#gonethttp)
    - [go:os](#goos)
    - [go:sync](#gosync)
    - [typego:memory](#typegomemory)
    - [typego:worker](#typegoworker)
  - [Concurrency](#concurrency)
  - [Memory Management](#memory-management)
- [Tooling](#tooling)
  - [CLI Reference](#cli-reference)
  - [Package Management](#package-management)
- [Ecosystem](#ecosystem)
  - [Performance](#performance)
  - [Runtime Comparison](#runtime-comparison)
- [Development](#development)
- [License](#license)

---

## Overview

TypeGo bridges the gap between Go's raw performance and TypeScript's developer experience. It is not just a JS runtime; it is a **hybrid runtime** where TypeScript code compiles JIT into the Go process, sharing memory and goroutines.

### Features

- **Direct Go Integration**: Import any Go package as a native TS module (`go:fmt`, `go:github.com/gin-gonic/gin`).
- **Standard Library Intrinsics**: Direct access to Go-native keywords and types like `go` routines, `makeChan`, `select`, `defer`, `ref`/`deref`, and `make`/`cap`.
- **Smarter Type Linker**: Automatic, recursive type generation for Go structs, interfaces, and methods. Supports struct embedding (`extends`) and nested type resolution.
- **True Parallelism**: Goroutine-based workers with zero-copy shared memory (`typego:memory`).
- **Modern Package Ecosystem**: Built-in CLI for managing Go dependencies with `typego.modules.json` and `typego.lock`.
- **Fast Developer Loop**: Hot-reloading dev server and ~0.2s interpreter startup. Compiles to single-binary with `--compile`.

---

## Quick Start

### Installation

```bash
go install github.com/repyh/typego/cmd/typego@latest
```

### Usage

**1. Initialize a Project**

```bash
typego init myapp
cd myapp
```

**2. Run Development Server**

```bash
# Watch mode with hot-reload
typego dev src/index.ts
```

**3. Build for Production**

```bash
# Compile to a standalone binary
typego build src/index.ts -o app
./app
```

---

## Language Reference

### Imports

TypeGo uses a unique import scheme to distinguish between TypeScript/JS modules and Go packages.

```typescript
// Import standard Go packages
import { Println } from "go:fmt";

// Import external Go modules (must be added via CLI first)
import { Default } from "go:github.com/gin-gonic/gin";

// Import TypeGo internal modules
import { Worker } from "typego:worker";

// Import relative TypeScript files
import { util } from "./util";
```

### Intrinsics

TypeGo exposes low-level Go primitives as global functions, allowing you to write "Go-like" TypeScript.

| Function | Type | Description |
|----------|--------|-------------|
| `go(fn, ...args)` | Concurrency | Launches a background goroutine. |
| `makeChan<T>(size?)` | Concurrency | Creates a synchronized Go channel. |
| `select(cases)` | Concurrency | Multiplexes channel operations (send/receive/default). |
| `ref(val)` | Memory | Creates a pointer handle (`Ref<T>`) to a value on the Go heap. |
| `deref(ptr)` | Memory | Dereferences a pointer or `Ref` object. |
| `make(Type, len, cap?)` | Memory | Allocates high-performance slices (TypedArrays). |
| `cap(v)` | Memory | Returns the capacity of a slice or buffer. |
| `copy(dst, src)` | Memory | Performs high-speed memory copying between buffers. |
| `sizeof(obj)` | Memory | Estimates the memory footprint of a JS/Go object. |
| `defer(fn)` | Control | Registers a function to run when the current `typego.scope` exits. |
| `panic(err)` | Control | Triggers a native Go panic. |
| `recover()` | Control | Recovers from a panic inside a `defer` block. |
| `iota` | Constant | Auto-incrementing compile-time constant. |

### Standard Library

TypeGo includes pre-bound versions of common Go standard library packages.

#### go:fmt
Print formatted output to stdout/stderr.

```typescript
import { Println, Sprintf } from "go:fmt";
Println("Hello", "World");
const msg = Sprintf("Value: %d", 42);
```

#### go:net/http
Make HTTP requests or run a server.

```typescript
import { Get, Post, ListenAndServe } from "go:net/http";

// Client
const resp = Get("https://example.com");
console.log(resp.Status);

// Server
ListenAndServe(":8080", (w, r) => {
    // Note: Handler signature adaptation might be needed depending on usage
});
```

#### go:os
File system interactions (sandboxed to CWD by default).

```typescript
import { ReadFile, WriteFile, Exit } from "go:os";

WriteFile("test.txt", "Hello TypeGo");
const content = ReadFile("test.txt");
Exit(0);
```

#### go:sync
Concurrency primitives.

```typescript
import { Spawn, Sleep } from "go:sync";

Spawn(async () => {
    await Sleep(100);
    console.log("Async work done");
});
```

#### typego:memory
Shared memory for workers.

```typescript
import { makeShared } from "typego:memory";
const buf = makeShared("myBuffer", 1024); // Globally accessible as 'myBuffer'
```

#### typego:worker
Thread-based worker spawning.

```typescript
import { Worker } from "typego:worker";
const w = new Worker("./worker.ts");
w.postMessage({ task: "compute" });
```

### Concurrency

TypeGo offers true parallelism via goroutines, distinct from Node.js's single-threaded event loop.

```typescript
import { Println } from "go:fmt";

// Channels + select
const ch = makeChan<number>(1);

go(() => {
  ch.send(42);
});

select([
  { chan: ch, recv: (v) => Println("received:", v) },
  { default: () => Println("no message") }
]);
```

### Memory Management

TypeGo manages memory automatically but provides tools for manual control when performance is critical.

- **`defer(fn)`**: Schedules cleanup execution (similar to Go's `defer`).
- **`ref(val)`**: Creates a pointer to a value, avoiding copy overhead for large structs.
- **Shared Memory**: Use `typego:memory` to share buffers between workers without serialization.

---

## Tooling

### CLI Reference

| Command | Description |
|---------|-------------|
| `typego run <file>` | Execute TypeScript (fast interpreter mode) |
| `typego dev <file>` | Development server with hot-reload |
| `typego build <file>` | Build standalone executable |
| `typego init [name]` | Scaffold new project (`--npm` for Node interop) |
| `typego types` | Generate `.d.ts` for Go imports |
| `typego add <pkg>` | Add a Go module dependency |
| `typego remove <pkg>` | Remove a Go module dependency |
| `typego list` | List configured Go dependencies |
| `typego update` | Update Go modules to latest versions |
| `typego outdated` | Check for newer Go module versions |
| `typego install` | Manually trigger JIT build/dependency resolution |
| `typego clean` | Reset build cache and temporary workspace |

### Package Management

TypeGo uses a `typego.modules.json` file to manage Go dependencies. This allows you to use any Go package in your TypeScript code.

```bash
# Add a Go package. Ecosystem will automatically resolve versions and sync types.
typego add github.com/gin-gonic/gin
```

TypeGo automatically manages a `.typego/` workspace, handling `go mod tidy`, JIT compilation, and TypeScript definition syncing behind the scenes.

---

## Ecosystem

### Performance

TypeGo is optimized for high-throughput I/O and true parallelism. For a detailed breakdown of the execution model, interop overhead, and optimization strategies, see **[OPTIMIZATION.md](OPTIMIZATION.md)**.

### Runtime Comparison

| Feature | TypeGo | Node.js | Deno | Bun |
|---------|:------:|:-------:|:----:|:---:|
| TypeScript native | ✅ | ⚠️ | ✅ | ✅ |
| True parallelism | ✅ | ❌ | ❌ | ❌ |
| Single binary | ✅ | ❌ | ✅ | ❌ |
| Shared memory | ✅ | ⚠️ | ⚠️ | ⚠️ |
| NPM ecosystem | ⚠️ | ✅ | ⚠️ | ✅ |

---

## Development

### Prerequisites

- Go 1.21+
- Node.js 18+ (for NPM packages)

### Building from Source

```bash
git clone https://github.com/repyh/typego.git
cd typego
go build -o typego.exe ./cmd/typego
```

### Running Examples

```bash
./typego run examples/01-hello-world.ts
./typego run examples/02-concurrency-basics.ts
./typego run examples/09-typego-stdlib.ts
```

---

## License

MIT
