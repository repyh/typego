package bridge

import (
	"fmt"
	"reflect"

	"github.com/dop251/goja"
)

// Binding represents a mapped Go struct in JS
type Binding struct {
	Name   string
	Target interface{}
}

// BindStruct uses reflection to expose all exported methods and fields of a struct to the JS VM
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
