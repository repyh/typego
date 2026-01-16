package bridge

import (
	"sync"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

type AsyncMutex struct {
	mu *sync.RWMutex
	el *eventloop.EventLoop
}

func NewAsyncMutex(mu *sync.RWMutex, el *eventloop.EventLoop) *AsyncMutex {
	return &AsyncMutex{mu: mu, el: el}
}

func (m *AsyncMutex) Lock(vm *goja.Runtime) goja.Value {
	p, resolve, _ := m.el.CreatePromise()
	go func() {
		m.mu.Lock()
		resolve(goja.Undefined())
	}()
	return p
}

func (m *AsyncMutex) Unlock() {
	m.mu.Unlock()
}

func (m *AsyncMutex) RLock(vm *goja.Runtime) goja.Value {
	p, resolve, _ := m.el.CreatePromise()
	go func() {
		m.mu.RLock()
		resolve(goja.Undefined())
	}()
	return p
}

func (m *AsyncMutex) RUnlock() {
	m.mu.RUnlock()
}

func BindMutex(vm *goja.Runtime, mu *sync.RWMutex, el *eventloop.EventLoop) goja.Value {
	am := NewAsyncMutex(mu, el)
	obj := vm.NewObject()
	_ = obj.Set("lock", am.Lock)
	_ = obj.Set("unlock", am.Unlock)
	_ = obj.Set("rlock", am.RLock)
	_ = obj.Set("runlock", am.RUnlock)
	return obj
}
