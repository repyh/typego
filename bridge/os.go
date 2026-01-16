package bridge

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

// OSModule implements the go/os package
type OSModule struct {
	Root string
}

// sanitizePath ensures the path is within the root directory
func (m *OSModule) sanitizePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(abs, m.Root) {
		return "", os.ErrPermission
	}

	return abs, nil
}

// WriteFile maps to os.WriteFile
func (m *OSModule) WriteFile(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		path := call.Argument(0).String()
		data := call.Argument(1).String()

		safePath, err := m.sanitizePath(path)
		if err != nil {
			panic(vm.NewTypeError("Sandbox violation or invalid path: ", err.Error()))
		}

		err = os.WriteFile(safePath, []byte(data), 0644)
		if err != nil {
			panic(vm.NewTypeError("os.WriteFile error: ", err.Error()))
		}

		return goja.Undefined()
	}
}

// RegisterOS injects the os functions into the runtime
func RegisterOS(vm *goja.Runtime) {
	wd, _ := os.Getwd()
	m := &OSModule{Root: wd}

	obj := vm.NewObject()
	obj.Set("WriteFile", m.WriteFile(vm))

	vm.Set("__go_os__", obj)
}
