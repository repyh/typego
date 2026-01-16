package bridge

import (
	"io"
	"net/http"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

// HTTPModule implements the go/net/http package
type HTTPModule struct{}

// Response wraps http.Response for JS
type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

// Get maps to http.Get
func (h *HTTPModule) Get(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		url := call.Argument(0).String()

		resp, err := http.Get(url)
		if err != nil {
			panic(vm.NewTypeError("http.Get error: ", err.Error()))
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		res := vm.NewObject()
		res.Set("Status", resp.Status)
		res.Set("StatusCode", resp.StatusCode)
		res.Set("Body", string(body))

		return res
	}
}

// RegisterHTTP injects the http functions into the runtime
func RegisterHTTP(vm *goja.Runtime, el *eventloop.EventLoop) {
	h := &HTTPModule{}

	obj := vm.NewObject()
	obj.Set("Get", h.Get(vm))

	obj.Set("Fetch", func(call goja.FunctionCall) goja.Value {
		url := call.Argument(0).String()
		p, resolve, reject := el.CreatePromise()

		go func() {
			resp, err := http.Get(url)
			if err != nil {
				reject(vm.NewTypeError("Fetch error: ", err.Error()))
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			el.RunOnLoop(func() {
				res := vm.NewObject()
				res.Set("Status", resp.Status)
				res.Set("StatusCode", resp.StatusCode)
				res.Set("Body", string(body))
				resolve(res)
			})
		}()

		return p
	})

	vm.Set("__go_http__", obj)
}
