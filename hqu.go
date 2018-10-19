package hqu

const (
	bucketSize  = 8
	maxFreelist = 64
)

var useFreelist = true

// Queuer interface for queue
type Queuer interface {
	Enqueue(v interface{})
	Dequeue() (v interface{}, ok bool)
}

// Stacker interface for stack
type Stacker interface {
	Push(v interface{})
	Pop() (v interface{}, ok bool)
}

type node struct {
	v interface{}

	prev, next *node
}

var (
	null = &node{nil, nil, nil}
)
