// Package bridge provides the JavaScript-to-Go binding layer for the TypeGo runtime.
//
// Bridge exposes Go functionality to the JavaScript environment running inside Goja.
// It provides modules that mirror Go's standard library (fmt, os, net/http) as well as
// TypeGo-specific features like shared memory, workers, and Node.js polyfills.
//
// # Module Registration
//
// Each module provides a Register* function that injects its functionality into
// the Goja runtime. These are typically called during engine initialization:
//
//	bridge.RegisterConsole(vm)  // console.log, console.error
//	bridge.RegisterFmt(vm)      // go/fmt.Println
//	bridge.RegisterOS(vm)       // go/os.WriteFile
//	bridge.RegisterHTTP(vm, el) // go/net/http.Get, Fetch
//	bridge.RegisterSync(vm, el) // Spawn, Sleep, Chan
//
// # Shared Memory
//
// Bridge provides two mechanisms for shared memory between the main thread and workers:
//
//   - MapSharedBuffer: Maps a Go byte slice to a JavaScript ArrayBuffer
//   - MemoryFactory: Creates named shared memory segments accessible by name
//
// Both mechanisms allow zero-copy access to the underlying memory from JavaScript.
//
// # Worker Support
//
// The WorkerHandle interface and RegisterWorker function enable spawning background
// workers that communicate via postMessage/onmessage, similar to Web Workers.
//
// # Node.js Compatibility
//
// The polyfills subpackage provides implementations of common Node.js globals
// (process, Buffer, setTimeout) for NPM package compatibility.
package bridge
