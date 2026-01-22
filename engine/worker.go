package engine

import (
	"fmt"

	"github.com/grafana/sobek"
	"github.com/repyh/typego/bridge/stdlib/worker"
	"github.com/repyh/typego/compiler"
)

type WorkerInstance struct {
	vm          *sobek.Runtime
	engine      *Engine
	inbox       chan interface{}
	stop        chan struct{}
	scriptPath  string
	onMessage   func(sobek.Value)
	autoRespawn bool
}

func (w *WorkerInstance) PostMessage(msg sobek.Value) {
	data := msg.Export()
	w.inbox <- data
}

func (w *WorkerInstance) Terminate() {
	w.autoRespawn = false
	close(w.stop)
}

func (e *Engine) SpawnWorker(scriptPath string, onMessage func(sobek.Value)) (worker.Handle, error) {
	inbox := make(chan interface{}, 100)
	stop := make(chan struct{})

	w := &WorkerInstance{
		scriptPath:  scriptPath,
		onMessage:   onMessage,
		inbox:       inbox,
		stop:        stop,
		autoRespawn: true,
	}

	e.startWorker(w)

	return w, nil
}

func (e *Engine) startWorker(w *WorkerInstance) {
	go func() {
		for {
			res, err := compiler.Compile(w.scriptPath, nil)
			if err != nil {
				fmt.Printf("Worker Compile Error [%s]: %v\n", w.scriptPath, err)
				return
			}

			workerEng := NewEngine(e.MemoryLimit, e.MemoryFactory)
			w.vm = workerEng.VM
			w.engine = workerEng

			worker.RegisterSelf(workerEng.VM, func(msg sobek.Value) {
				data := msg.Export()
				e.EventLoop.RunOnLoop(func() {
					val := e.VM.ToValue(data)
					w.onMessage(val)
				})
			})

			// Run Loop
			go func() {
				_, err := workerEng.Run(res.JS)
				if err != nil {
					fmt.Printf("Worker Runtime Error [%s]: %v\n", w.scriptPath, err)
				}
				workerEng.EventLoop.WGAdd(1)
				workerEng.EventLoop.Start()
			}()

			// Message Bridge
			bridgeDone := make(chan struct{})
			go func() {
				defer close(bridgeDone)
				for {
					select {
					case msg := <-w.inbox:
						workerEng.EventLoop.RunOnLoop(func() {
							val := workerEng.VM.ToValue(msg)
							if onMsg := workerEng.VM.GlobalObject().Get("onmessage"); onMsg != nil {
								if fn, ok := sobek.AssertFunction(onMsg); ok {
									event := workerEng.VM.NewObject()
									_ = event.Set("data", val)
									_, _ = fn(workerEng.VM.GlobalObject(), event)
								}
							}
						})
					case <-w.stop:
						workerEng.EventLoop.Stop()
						return
					case <-workerEng.Context().Done():
						return
					}
				}
			}()

			<-bridgeDone
			if !w.autoRespawn {
				return
			}
			fmt.Printf("Worker [%s] exited unexpectedly, respawning...\n", w.scriptPath)
		}
	}()
}
