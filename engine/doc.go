// Package engine provides the core JavaScript execution environment for TypeGo.
//
// Engine wraps the Goja JavaScript runtime with an event loop, memory management,
// and bridge bindings to Go packages. It serves as the main entry point for
// running TypeScript/JavaScript code in the TypeGo runtime.
//
// # Creating an Engine
//
// Use NewEngine to create a fully initialized runtime with all standard bindings:
//
//	eng := engine.NewEngine(128*1024*1024, nil) // 128MB memory limit
//	defer eng.EventLoop.Stop()
//
//	eng.EventLoop.RunOnLoop(func() {
//	    val, err := eng.Run(`console.log("Hello from TypeGo")`)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	})
//
//	eng.EventLoop.Start()
//
// # Memory Management
//
// The memoryLimit parameter sets a soft limit on JavaScript heap size. The engine
// monitors memory usage and can trigger emergency cleanup if limits are exceeded.
//
// # Workers
//
// The engine supports spawning worker threads via the SpawnWorker method. Workers
// run in isolated Goja runtimes but can share memory through the MemoryFactory.
//
// # Event Loop
//
// All JavaScript execution must occur on the event loop. Use RunOnLoop to schedule
// work, and WGAdd/WGDone to track pending async operations. The loop automatically
// stops when all work is complete.
package engine
