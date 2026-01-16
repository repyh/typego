package bridge

import (
	"github.com/dop251/goja"
)

// ToArrayBuffer converts a Go byte slice to a JS ArrayBuffer
func ToArrayBuffer(vm *goja.Runtime, data []byte) goja.Value {
	// Default Goja behavior creates a copy
	return vm.ToValue(vm.NewArrayBuffer(data))
}

// MapSharedBuffer exposes a Go slice as a JS TypedArray
func MapSharedBuffer(vm *goja.Runtime, name string, data []byte) {
	buf := vm.NewArrayBuffer(data)
	// Create a view (Uint8Array)
	view := vm.ToValue(vm.Get("Uint8Array")).ToObject(vm)

	// Create new instance: new Uint8Array(buffer)
	typedArray, _ := vm.New(view, vm.ToValue(buf))

	_ = vm.GlobalObject().Set(name, typedArray)
}
