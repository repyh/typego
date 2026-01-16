package bridge

import (
	"fmt"

	"github.com/dop251/goja"
)

// FmtModule implements the go/fmt package
type FmtModule struct{}

// Println maps to fmt.Println
func (f *FmtModule) Println(call goja.FunctionCall) goja.Value {
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println(args...)
	return goja.Undefined()
}

// RegisterFmt injects the fmt functions into the runtime
func RegisterFmt(vm *goja.Runtime) {
	f := &FmtModule{}

	// Create a namespace object for fmt
	obj := vm.NewObject()
	obj.Set("Println", f.Println)

	// Store it globally for the internal loader to find
	vm.Set("__go_fmt__", obj)
}
