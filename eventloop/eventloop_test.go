package eventloop

import (
	"testing"
	"time"

	"github.com/grafana/sobek"
)

func TestCreatePromiseDoubleResolve(t *testing.T) {
	vm := sobek.New()
	el := NewEventLoop(vm)
	el.SetAutoStop(false)

	go el.Start()
	defer el.Stop()

	_, resolve, _ := el.CreatePromise()

	// Call resolve twice. This should NOT panic if fixed.
	// Without fix, it causes negative WaitGroup counter.
	resolve(nil)
	resolve(nil)

	// Give it some time to process tasks
	time.Sleep(50 * time.Millisecond)
}
