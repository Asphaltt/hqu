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

func TestStack0(t *testing.T) {
	s := &Stack{}

	for i := 0; i < bucketSize+1; i++ {
		s.Push(i)
	}

	if len(s.buckets) != 2 {
		t.Logf("length of buckets is not 2, is %d", len(s.buckets))
		t.Fail()
	}

	s.Pop()

	if len(s.buckets) != 1 {
		t.Logf("length of buckets is not 1, is %d", len(s.buckets))
		t.Fail()
	}

	if len(s.freelist) != 1 {
		t.Logf("length of freelist is not 1, is %d", len(s.freelist))
		t.Fail()
	}
}

func TestStack1(t *testing.T) {
	s := &Stack{}

	for i := 0; i < bucketSize*maxFreelist+1; i++ {
		s.Push(i)
	}

	if len(s.buckets) != maxFreelist+1 {
		t.Logf("length of buckets is not %d, is %d", maxFreelist+1, len(s.buckets))
		t.Fail()
	}

	for i := 0; i < bucketSize*(maxFreelist*3/4)+1; i++ {
		s.Pop()
	}

	if len(s.buckets) != maxFreelist/4 {
		t.Logf("length of buckets is not %d, is %d", maxFreelist/4, len(s.buckets))
		t.Fail()
	}

	if cap(s.orgBuckets) != maxFreelist/2 {
		t.Logf("capacity of under array is not %d, is %d", maxFreelist/2, cap(s.orgBuckets))
		t.Fail()
	}
}

func TestStackSize(t *testing.T) {
	s := &Stack{}

	for i := 0; i < bucketSize+1; i++ {
		s.Push(i)
	}

	if s.Size() != bucketSize+1 {
		t.Logf("size of stack is not %d, is %d", bucketSize+1, s.Size())
		t.Fail()
	}

	for i := 0; i < bucketSize; i++ {
		s.Pop()
	}

	if s.Size() != 1 {
		t.Logf("size of stack is not %d, is %d", 1, s.Size())
		t.Fail()
	}
}
