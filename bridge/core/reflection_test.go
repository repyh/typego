package core

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/grafana/sobek"
)

type BenchmarkStruct struct {
	Map map[string]int
}

func BenchmarkBindStruct_Map(b *testing.B) {
	vm := sobek.New()
	data := BenchmarkStruct{
		Map: make(map[string]int),
	}
	for i := 0; i < 100; i++ {
		data.Map[string(rune(i))] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We are testing bindStruct -> bindMap path
		// We can call bindValue directly to avoid overhead of creating globals
		_, _ = bindValue(vm, reflect.ValueOf(data), make(map[uintptr]sobek.Value))
	}
}

type StringerInt int

func (s StringerInt) String() string {
	return "S-" + strconv.Itoa(int(s))
}

func TestBindMap_StringerKey(t *testing.T) {
	vm := sobek.New()
	data := make(map[StringerInt]int)
	data[StringerInt(1)] = 100

	val, err := bindValue(vm, reflect.ValueOf(data), make(map[uintptr]sobek.Value))
	if err != nil {
		t.Fatal(err)
	}

	obj := val.(*sobek.Object)
	if obj.Get("S-1").Export() != int64(100) {
		t.Errorf("Expected key 'S-1' to have value 100, got %v. Keys: %v", obj.Get("S-1"), obj.Keys())
	}
}

func BenchmarkBindMap_IntKey(b *testing.B) {
	vm := sobek.New()
	data := make(map[int]int)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bindValue(vm, reflect.ValueOf(data), make(map[uintptr]sobek.Value))
	}
}
