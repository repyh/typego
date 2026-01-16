// Package polyfills provides Node.js-compatible globals for the Goja runtime.
//
// Many NPM packages assume a Node.js environment with globals like process,
// Buffer, and setTimeout. This package provides lightweight polyfills that
// enable those packages to work in the TypeGo runtime.
//
// # Available Polyfills
//
//   - process: Environment variables, platform info, cwd, argv
//   - Buffer: from(), alloc() for basic buffer operations
//   - Timers: setTimeout, setInterval, clearInterval
//
// # Usage
//
// Call EnableAll to inject all polyfills at once:
//
//	polyfills.EnableAll(vm, eventLoop)
//
// Or enable individual polyfills as needed:
//
//	polyfills.EnableProcess(vm)
//	polyfills.EnableBuffer(vm)
//	polyfills.EnableTimers(vm, eventLoop)
//
// # Process Object
//
// The process polyfill provides:
//
//   - process.env: Maps to Go's os.Environ()
//   - process.platform: Returns runtime.GOOS
//   - process.cwd(): Returns os.Getwd()
//   - process.argv: Returns os.Args
//   - process.version: Returns runtime.Version()
//
// The FORCE_COLOR environment variable is automatically set to enable
// colored output for libraries like chalk.
//
// # Buffer Object
//
// A minimal Buffer implementation using Uint8Array:
//
//   - Buffer.from(string): Converts string to Uint8Array
//   - Buffer.alloc(size): Creates a zero-filled Uint8Array
//
// # Timers
//
// Timer functions integrate with the event loop to properly track async work:
//
//   - setTimeout(fn, ms): Schedules a one-time callback
//   - setInterval(fn, ms): Schedules a recurring callback
//   - clearInterval(id): Stops a recurring interval
package polyfills
