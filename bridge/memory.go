package bridge

import (
	"runtime"

	"github.com/dop251/goja"
)

// MemoryModule provides access to Go's memory statistics
type MemoryModule struct{}

// GetStats returns memory statistics to JS
func (m *MemoryModule) GetStats(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)

		obj := vm.NewObject()
		obj.Set("alloc", ms.Alloc)
		obj.Set("totalAlloc", ms.TotalAlloc)
		obj.Set("sys", ms.Sys)
		obj.Set("numGC", ms.NumGC)

		return obj
	}
}

// RegisterMemory injects the memory object and Ptr constructor into the runtime
func RegisterMemory(vm *goja.Runtime) {
	m := &MemoryModule{}
	vm.Set("goMemory", m.GetStats(vm))

	// Ptr constructor: allows wrapping a value to pass by reference in JS
	vm.Set("Ptr", func(call goja.ConstructorCall) *goja.Object {
		obj := vm.NewObject()
		val := call.Argument(0)
		_ = obj.Set("value", val)
		return obj
	})
}
