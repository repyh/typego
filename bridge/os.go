package bridge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

// OSModule implements the go/os package
type OSModule struct {
	Root string
}

// sanitizePath ensures the path is within the root directory and resolves symlinks
func (m *OSModule) sanitizePath(path string) (string, error) {
	// 1. Resolve to absolute path
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// 2. Resolve symlinks (The "Anti-Escape")
	// This prevents pointing a symlink inside the jail to a file outside
	realPath, err := filepath.EvalSymlinks(abs)
	if err != nil {
		// If path doesn't exist, we still check the parent
		if os.IsNotExist(err) {
			parentDir := filepath.Dir(abs)
			realParent, pErr := filepath.EvalSymlinks(parentDir)
			if pErr == nil {
				realPath = filepath.Join(realParent, filepath.Base(abs))
			} else {
				return "", pErr
			}
		} else {
			return "", err
		}
	}

	// 3. Normalized Rel check
	rel, err := filepath.Rel(m.Root, realPath)
	if err != nil {
		return "", err
	}

	// If the path starts with ".." it is outside the root
	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", os.ErrPermission
	}

	return realPath, nil
}

// WriteFile maps to os.WriteFile
func (m *OSModule) WriteFile(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		path := call.Argument(0).String()
		data := call.Argument(1).String()

		safePath, err := m.sanitizePath(path)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("sandbox violation: %v", err)))
		}

		err = os.WriteFile(safePath, []byte(data), 0644)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("os.WriteFile error: %v", err)))
		}

		return goja.Undefined()
	}
}

// ReadFile maps to os.ReadFile
func (m *OSModule) ReadFile(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		path := call.Argument(0).String()

		safePath, err := m.sanitizePath(path)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("sandbox violation: %v", err)))
		}

		data, err := os.ReadFile(safePath)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("os.ReadFile error: %v", err)))
		}

		return vm.ToValue(string(data))
	}
}

// RegisterOS injects the os functions into the runtime
func RegisterOS(vm *goja.Runtime) {
	wd, _ := os.Getwd()
	absRoot, _ := filepath.Abs(wd)
	m := &OSModule{Root: absRoot}

	obj := vm.NewObject()
	obj.Set("WriteFile", m.WriteFile(vm))
	obj.Set("ReadFile", m.ReadFile(vm))

	vm.Set("__go_os__", obj)
}
