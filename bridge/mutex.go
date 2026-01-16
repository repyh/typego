package bridge

import (
	"sync"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

// AsyncMutex wraps a sync.RWMutex with async-friendly locking operations.
// Lock and RLock return Promises instead of blocking, allowing the JavaScript
// event loop to continue processing while waiting for the lock.
type AsyncMutex struct {
	mu *sync.RWMutex
	el *eventloop.EventLoop
}

// NewAsyncMutex creates an AsyncMutex wrapping the given RWMutex.
// The EventLoop is used to create Promises for async lock acquisition.
func NewAsyncMutex(mu *sync.RWMutex, el *eventloop.EventLoop) *AsyncMutex {
	return &AsyncMutex{mu: mu, el: el}
}

// Lock acquires an exclusive write lock asynchronously.
// Returns a Promise that resolves when the lock is acquired.
func (m *AsyncMutex) Lock(vm *goja.Runtime) goja.Value {
	p, resolve, _ := m.el.CreatePromise()
	go func() {
		m.mu.Lock()
		resolve(goja.Undefined())
	}()
	return p
}

// Unlock releases a previously acquired write lock.
func (m *AsyncMutex) Unlock() {
	m.mu.Unlock()
}

// RLock acquires a shared read lock asynchronously.
// Multiple readers can hold the lock simultaneously.
// Returns a Promise that resolves when the lock is acquired.
func (m *AsyncMutex) RLock(vm *goja.Runtime) goja.Value {
	p, resolve, _ := m.el.CreatePromise()
	go func() {
		m.mu.RLock()
		resolve(goja.Undefined())
	}()
	return p
}

// RUnlock releases a previously acquired read lock.
func (m *AsyncMutex) RUnlock() {
	m.mu.RUnlock()
}

// BindMutex creates a JavaScript object exposing mutex operations.
// The returned object has lock(), unlock(), rlock(), and runlock() methods.
// Lock and rlock are async (return Promises), while unlock and runlock are sync.
//
// Example usage in JavaScript:
//
//	await mutex.lock();
//	try {
//	    // critical section
//	} finally {
//	    mutex.unlock();
//	}
func BindMutex(vm *goja.Runtime, mu *sync.RWMutex, el *eventloop.EventLoop) goja.Value {
	am := NewAsyncMutex(mu, el)
	obj := vm.NewObject()
	_ = obj.Set("lock", am.Lock)
	_ = obj.Set("unlock", am.Unlock)
	_ = obj.Set("rlock", am.RLock)
	_ = obj.Set("runlock", am.RUnlock)
	return obj
}
