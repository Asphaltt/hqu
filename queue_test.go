package hqu

import "testing"

func TestQueueEnqueue(t *testing.T) {
	hq := &Queue{}

	hq.Enqueue(1)

	v, ok := hq.Dequeue()
	if !ok || v != 1 {
		t.Fail()
	}
}

func TestQueueDequeue(t *testing.T) {
	hq := &Queue{}

	v, ok := hq.Dequeue()
	if ok || v != nil {
		t.Fail()
	}
}
