package hqu

import (
	"sync"
)

// Stack is a high memory efficient stack.
// It uses bucket to cache data and reuses the buckets.
type Stack struct {
	sync.Mutex

	pos int // the position where a new value should be inserted

	buckets  [][]interface{}
	freelist [][]interface{}
}

// Push pushes a value into the stack.
func (s *Stack) Push(v interface{}) {
	s.Lock()

	bp := s.pos / bucketSize

	// look up bucket
	var bkt []interface{}
	if bp == len(s.buckets) {
		if useFreelist && len(s.freelist) != 0 {
			// reuse bucket
			idx := len(s.freelist) - 1
			bkt = s.freelist[idx]
			s.freelist[idx] = nil
			s.freelist = s.freelist[:idx]
		} else {
			bkt = make([]interface{}, bucketSize) // create bucket
		}
		s.buckets = append(s.buckets, bkt)
	} else {
		bkt = s.buckets[bp]
	}
	bkt[s.pos%bucketSize] = v
	s.pos++

	s.Unlock()
}

// Pop pops a value from the stack.
func (s *Stack) Pop() (v interface{}, ok bool) {
	s.Lock()
	if s.pos == 0 {
		s.Unlock()
		return nil, false
	}
	s.pos--

	bp, qp := s.pos/bucketSize, s.pos%bucketSize

	// lookup bucket
	bkt := s.buckets[bp]
	v, ok = bkt[qp], true
	bkt[qp] = nil // free the value

	if qp == 0 {
		s.buckets[bp] = nil // free the bucket
		s.buckets = s.buckets[:len(s.buckets)-1]

		// reuse bucket
		if useFreelist && len(s.freelist) < maxFreelist {
			s.freelist = append(s.freelist, bkt)
		}
	}

	s.Unlock()
	return
}
