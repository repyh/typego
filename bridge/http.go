package bridge

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

// Default HTTP client with production-ready timeouts
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

const maxResponseBodySize = 50 * 1024 * 1024 // 50MB

// HTTPModule implements the go/net/http package
type HTTPModule struct{}

// Get maps to http.Get
func (h *HTTPModule) Get(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		url := call.Argument(0).String()

		resp, err := httpClient.Get(url)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("http.Get error: %v", err)))
		}
		defer resp.Body.Close()

		// Safeguard against OOM by limiting the reader
		limitedReader := &io.LimitedReader{R: resp.Body, N: maxResponseBodySize}
		body, err := io.ReadAll(limitedReader)
		if err != nil && err != io.EOF {
			panic(vm.NewTypeError(fmt.Sprintf("http body read error: %v", err)))
		}

		if limitedReader.N <= 0 {
			panic(vm.NewTypeError(fmt.Sprintf("http response too large (max %d MB)", maxResponseBodySize/1024/1024)))
		}

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
			resp, err := httpClient.Get(url)
			if err != nil {
				el.RunOnLoop(func() {
					reject(vm.NewTypeError(fmt.Sprintf("Fetch error: %v", err)))
				})
				return
			}
			defer resp.Body.Close()

			limitedReader := &io.LimitedReader{R: resp.Body, N: maxResponseBodySize}
			body, err := io.ReadAll(limitedReader)

			el.RunOnLoop(func() {
				if err != nil && err != io.EOF {
					reject(vm.NewTypeError(fmt.Sprintf("Fetch body read error: %v", err)))
					return
				}
				if limitedReader.N <= 0 {
					reject(vm.NewTypeError(fmt.Sprintf("Fetch response too large (max %d MB)", maxResponseBodySize/1024/1024)))
					return
				}

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
