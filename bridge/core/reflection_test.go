package core

import (
	"testing"

	"github.com/grafana/sobek"
)

func BenchmarkBindMap(b *testing.B) {
	vm := sobek.New()
	data := make(map[string]int)
	for i := 0; i < 100; i++ {
		data[string(rune(i))] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BindStruct(vm, "data", data)
	}
}

func BenchmarkBindMapIntKeys(b *testing.B) {
	vm := sobek.New()
	data := make(map[int]int)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BindStruct(vm, "data", data)
	}
}

func TestBindMap(t *testing.T) {
	vm := sobek.New()
	data := map[string]int{"a": 1, "b": 2}

	err := BindStruct(vm, "data", data)
	if err != nil {
		t.Fatalf("BindStruct failed: %v", err)
	}

	val, _ := vm.RunString("data.a")
	if val.Export() != int64(1) {
		t.Errorf("Expected 1, got %v", val.Export())
	}
}

func TestBindMapIntKeys(t *testing.T) {
	vm := sobek.New()
	data := map[int]int{1: 10, 2: 20}

	err := BindStruct(vm, "data", data)
	if err != nil {
		t.Fatalf("BindStruct failed: %v", err)
	}

	val, _ := vm.RunString("data['1']")
	if val.Export() != int64(10) {
		t.Errorf("Expected 10, got %v", val.Export())
	}
}
