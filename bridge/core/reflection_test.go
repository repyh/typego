package core

import (
	"testing"

	"github.com/grafana/sobek"
)

func BenchmarkBindMap(b *testing.B) {
	vm := sobek.New()
	m := make(map[string]int)
	for i := 0; i < 1000; i++ {
		m[string(rune(i))] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BindStruct(vm, "testMap", m)
	}
}

func BenchmarkBindMapIntKeys(b *testing.B) {
	vm := sobek.New()
	m := make(map[int]int)
	for i := 0; i < 1000; i++ {
		m[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BindStruct(vm, "testMapInt", m)
	}
}
