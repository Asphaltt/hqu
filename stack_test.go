package hqu

import (
	"testing"
)

func TestStackPush(t *testing.T) {
	hq := &Stack{}

	hq.Push("test")
	v, ok := hq.Pop()
	if !ok || v != "test" {
		t.Fail()
	}
}

func TestStackPop(t *testing.T) {
	hq := &Stack{}

	v, ok := hq.Pop()
	if ok || v != nil {
		t.Fail()
	}
}

func BenchmarkStackWithFreelist(b *testing.B) {
	useFreelist = true
	hq := &Stack{}
	for i := 0; i < b.N; i++ {
		hq.Push(i)
	}

	for i := 0; i < b.N; i++ {
		hq.Pop()
	}
}

func BenchmarkStackWithoutFreelist(b *testing.B) {
	useFreelist = false
	hq := &Stack{}
	for i := 0; i < b.N; i++ {
		hq.Push(i)
	}

	for i := 0; i < b.N; i++ {
		hq.Pop()
	}
}
