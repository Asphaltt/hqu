package hqu

import (
	"sync"
)

// Queue is a high memory efficient queue.
// It uses bucket to cache data and reuses the buckets.
type Queue struct {
	sync.Mutex

	// front is the location where a value should be dequeued
	// rear is the location where a new value should be enqueued
	front, rear int

	buckets  [][]interface{}
	freelist [][]interface{}

	// for reallocation of buckets
	orgBuckets [][]interface{}
}

// Size gets the count of elements in queue
func (q *Queue) Size() int {
	q.Lock()
	size := q.rear - q.front
	q.Unlock()
	return size
}

// Enqueue pushes a value into the queue.
func (q *Queue) Enqueue(v interface{}) {
	q.Lock()

	bp := q.rear / bucketSize

	var bkt []interface{}
	if bp == len(q.buckets) {
		if useFreelist && len(q.freelist) != 0 {
			// reuse bucket
			idx := len(q.freelist) - 1
			bkt = q.freelist[idx]
			q.freelist[idx] = nil
			q.freelist = q.freelist[:idx]
		} else {
			// create bucket
			bkt = make([]interface{}, bucketSize)
		}
		// realloc freelist automatically when necessary, TestSlice in slice_test.go
		q.buckets = append(q.buckets, bkt)
		q.orgBuckets = q.buckets
	} else {
		bkt = q.buckets[bp]
	}

	bkt[q.rear%bucketSize] = v
	q.rear++

	q.Unlock()
}

// Dequeue pops a value from the queue.
func (q *Queue) Dequeue() (v interface{}, ok bool) {
	q.Lock()
	if q.rear == q.front {
		q.Unlock()
		return nil, false
	}

	bp, fp := q.front/bucketSize, q.front%bucketSize

	bkt := q.buckets[bp]
	v, ok = bkt[fp], true
	bkt[fp] = nil // free the value

	q.front++
	if q.front%bucketSize == 0 {
		q.buckets[bp] = nil // free the bucket
		q.buckets = q.buckets[1:]

		// recude memory usage when no reallocation(append in Enqueue), TestSlice1 in slice_test.go
		if len(q.buckets) <= cap(q.orgBuckets)/4 { // the usage is less than or equal to a quater of the capacity
			tmp := make([][]interface{}, cap(q.orgBuckets)/2)
			n := copy(tmp, q.buckets)
			q.buckets = tmp[:n]
			q.orgBuckets = tmp
		}

		// reuse bucket
		if useFreelist && len(q.freelist) < maxFreelist {
			q.freelist = append(q.freelist, bkt)
		}

		q.front = 0
		q.rear -= bucketSize
	}

	q.Unlock()
	return
}
