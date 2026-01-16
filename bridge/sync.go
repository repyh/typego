package bridge

import (
	"time"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

// SyncModule implements the go/sync package
type SyncModule struct {
	el *eventloop.EventLoop
}

func (m *SyncModule) Spawn(call goja.FunctionCall) goja.Value {
	fn, ok := goja.AssertFunction(call.Argument(0))
	if !ok {
		panic(m.el.VM.NewTypeError("Spawn expects a function"))
	}

	go func() {
		m.el.RunOnLoop(func() {
			val, err := fn(goja.Undefined())
			if err != nil {
				return
			}

			// Handle async functions by waiting for the returned promise
			if obj := val.ToObject(m.el.VM); obj != nil {
				then := obj.Get("then")
				if then != nil && !goja.IsUndefined(then) {
					if thenFn, ok := goja.AssertFunction(then); ok {
						m.el.WGAdd(1)
						done := m.el.VM.ToValue(func(goja.FunctionCall) goja.Value {
							m.el.WGDone()
							return goja.Undefined()
						})
						_, _ = thenFn(val, done, done)
					}
				}
			}
		})
	}()

	return goja.Undefined()
}

func (m *SyncModule) Sleep(call goja.FunctionCall) goja.Value {
	ms := call.Argument(0).ToInteger()
	p, resolve, _ := m.el.CreatePromise()

	go func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		m.el.RunOnLoop(func() {
			resolve(goja.Undefined())
		})
	}()

	return p
}

func RegisterSync(vm *goja.Runtime, el *eventloop.EventLoop) {
	m := &SyncModule{el: el}
	obj := vm.NewObject()
	obj.Set("Spawn", m.Spawn)
	obj.Set("Sleep", m.Sleep)

	vm.Set("Chan", func(call goja.ConstructorCall) *goja.Object {
		ch := make(chan goja.Value, 100)
		res := vm.NewObject()
		res.Set("send", func(c goja.FunctionCall) goja.Value {
			ch <- c.Argument(0)
			return goja.Undefined()
		})
		res.Set("recv", func(c goja.FunctionCall) goja.Value {
			p, resolve, _ := el.CreatePromise()
			go func() { resolve(<-ch) }()
			return p
		})
		return res
	})

	vm.Set("__go_sync__", obj)
}
