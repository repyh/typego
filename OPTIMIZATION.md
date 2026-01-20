# TypeGo Execution Model & Performance Optimization

This document outlines the architectural differences between the embedded JavaScript environment (Goja) and the TypeGo native bridge, providing technical guidelines for optimal implementation.

## Execution Contexts

TypeGo applications operate across two distinct execution contexts:

1.  **Script Context (Goja VM)**
    *   **Type:** Interpreted JavaScript (ES5/ES6 compliant)
    *   **Characteristics:** Single-threaded, managed heap, dynamic typing.
    *   **Use Case:** Business logic, flow control, JSON manipulation, string processing.

2.  **Native Context (Go Runtime)**
    *   **Type:** Compiled Machine Code (Go)
    *   **Characteristics:** Multi-threaded (Goroutines), strictly typed, manual memory management (via arenas/shared buffers).
    *   **Use Case:** I/O operations, CPU-bound concurrent tasks, system calls.

## Interop Overhead

Every function call traversing the JS-to-Go boundary incurs a context switching overhead. Use the following metrics for architectural decisions:

*   **Bridge Call Latency:** ~1µs - 5µs per call.
*   **Data Marshaling:** Auto-conversion of complex types (e.g., deeply nested objects) incurs reflection costs.

## Functional Domain Analysis

### 1. I/O & Networking

| Operation | Implementation | Context | Technical Justification |
| :--- | :--- | :---: | :--- |
| **HTTP Requests** | `go:net/http` | ✅ Native | Go's `net` package utilizes the OS-level non-blocking I/O poller (epoll/kqueue), offering vastly superior throughput and connection handling compared to a purely interpreted implementation. |
| **File I/O** | `go:os` | ✅ Native | Direct syscall usage with minimal GC pressure. Provides strict path validation and isolation not native to the JS VM. |

### 2. Concurrency & Parallelism

| Mechanism | Implementation | Type | Characteristics |
| :--- | :--- | :---: | :--- |
| **JS Async/Await** | `Promise` | ❌ Cooperative | Single-threaded. Suitable for effectively hiding I/O latency within the script context but cannot parallelize CPU load. |
| **TypeGo Spawn** | `go:sync.Spawn` | ✅ Preemptive | Spawns a dedicated OS thread (managed by Go scheduler). Suitable for CPU-bound tasks that would otherwise block the JS event loop. |

### 3. Data Processing

| Task | Preferred Context | Rationale |
| :--- | :---: | :--- |
| **JSON Parsing** | ✅ Script | Goja's native JSON parser is highly optimized for dynamic object creation. Go structural unmarshaling incurs reflection overhead. |
| **RegEx** | ✅ Script | Faster for standard validation. Go's `regexp` package guarantees linear time complexity $O(n)$ which is safer but slower than V8/Goja backtracking engines for simple patterns. |
| **Binary Manipulation** | ✅ Native | Use `SharedArrayBuffer` / `makeShared`. Go provides efficient slice manipulation without V8 object overhead. |

## Optimization Strategies

### A. Batch Bridge Calls
Minimize the frequency of crossing the bridge boundary.

*   **Inefficient:** Calling a Go function inside a tight JS loop.
*   **Optimized:** Passing an array to Go and iterating within the Native Context.

### B. Memory Management
*   **Script Context:** Avoid creating excessive temporary objects to reduce GC pressure/pauses.
*   **Shared Memory:** Use `typego:memory` for large data sets (> 10KB) to avoid serialization costs during transfer.

### API Reference Map

| Capability | Script Standard (JS) | Native Bridge (Go) | Recommendation |
| :--- | :--- | :--- | :--- |
| **Debug Logging** | `console.log` | `fmt.Println` | ✅ `console.log` |
| **Prod Logging** | - | `fmt` | ✅ `fmt` |
| **Async Wait** | `setTimeout` | - | ✅ `setTimeout` |
| **Blocking Wait** | - | `time.Sleep` | ✅ `time.Sleep` |
| **Error Handling** | `try/catch` | `error` | ✅ `try/catch` |
