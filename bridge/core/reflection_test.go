package core

import (
	"fmt"
	"testing"

	"github.com/grafana/sobek"
)

func BenchmarkBindMapStringKeys(b *testing.B) {
	vm := sobek.New()
	data := make(map[string]int)
	for i := 0; i < 1000; i++ {
		data[fmt.Sprintf("key-%d", i)] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BindStruct(vm, "benchMap", data)
	}
}

func BenchmarkBindMapIntKeys(b *testing.B) {
	vm := sobek.New()
	data := make(map[int]int)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BindStruct(vm, "benchMapInt", data)
	}
}
