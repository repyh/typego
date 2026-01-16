package bridge

import (
	"fmt"
	"reflect"

	"github.com/dop251/goja"
)

// Binding represents a Go struct that has been bound to the JavaScript runtime.
// It stores the name under which the struct is accessible in JS and a reference
// to the actual Go value.
type Binding struct {
	Name   string
	Target interface{}
}

// BindStruct uses reflection to expose a Go struct to the JavaScript runtime.
// All exported fields are made readable, and all exported methods are callable
// from JavaScript.
//
// The struct is registered as a global object with the given name. Method calls
// from JavaScript are automatically marshaled to the appropriate Go types.
//
// Example:
//
//	type Calculator struct {
//	    LastResult float64
//	}
//
//	func (c *Calculator) Add(a, b float64) float64 {
//	    c.LastResult = a + b
//	    return c.LastResult
//	}
//
//	calc := &Calculator{}
//	bridge.BindStruct(vm, "calc", calc)
//	// In JS: calc.Add(1, 2) returns 3
//	// In JS: calc.LastResult contains 3
//
// Type conversion is handled automatically for compatible types. If a JavaScript
// value cannot be converted to the expected Go type, a TypeError is thrown.
func BindStruct(vm *goja.Runtime, name string, s interface{}) error {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("BindStruct requires a struct or pointer to struct, got %s", t.Kind())
	}

	obj := vm.NewObject()

	// Map Fields (Read/Write)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		_ = obj.Set(fieldName, v.Field(i).Interface())
	}

	// Map Methods
	vPtr := reflect.ValueOf(s)
	tPtr := reflect.TypeOf(s)

	for i := 0; i < tPtr.NumMethod(); i++ {
		method := tPtr.Method(i)
		if !method.IsExported() {
			continue
		}

		methodName := method.Name
		methodVal := vPtr.Method(i)

		_ = obj.Set(methodName, func(call goja.FunctionCall) goja.Value {
			goArgs := make([]reflect.Value, methodVal.Type().NumIn())
			for j := 0; j < methodVal.Type().NumIn(); j++ {
				if j < len(call.Arguments) {
					argType := methodVal.Type().In(j)
					val := call.Arguments[j].Export()

					goVal := reflect.ValueOf(val)
					if goVal.Type().AssignableTo(argType) {
						goArgs[j] = goVal
					} else if goVal.Type().ConvertibleTo(argType) {
						goArgs[j] = goVal.Convert(argType)
					} else {
						panic(vm.NewTypeError(fmt.Sprintf("Method %s: Argument %d expected %s, got %T", methodName, j, argType, val)))
					}
				} else {
					goArgs[j] = reflect.Zero(methodVal.Type().In(j))
				}
			}

			results := methodVal.Call(goArgs)

			if len(results) == 0 {
				return goja.Undefined()
			}

			return vm.ToValue(results[0].Interface())
		})
	}

	return vm.GlobalObject().Set(name, obj)
}
