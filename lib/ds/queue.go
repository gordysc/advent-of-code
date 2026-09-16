package ds

// Queue is a first-in, first-out collection. Items are appended to a slice and
// a head index moves forward on Pop; the slice is compacted once the dead
// prefix grows large, so Pop is amortised O(1) without a linked list.
type Queue[T any] struct {
	items []T
	head  int
}

// NewQueue returns a queue holding the given items in order.
func NewQueue[T any](items ...T) *Queue[T] {
	q := &Queue[T]{}
	q.Push(items...)

	return q
}

// Push adds items to the back of the queue.
func (q *Queue[T]) Push(items ...T) {
	q.items = append(q.items, items...)
}

// Pop removes and returns the front item. It panics when the queue is empty;
// check [Queue.Len] first.
func (q *Queue[T]) Pop() T {
	item := q.items[q.head]

	// Clear the slot so the garbage collector can free what it pointed to.
	var zero T
	q.items[q.head] = zero
	q.head++

	// Compact when more than half the backing slice is dead space.
	if q.head > 32 && q.head*2 > len(q.items) {
		q.items = append(q.items[:0], q.items[q.head:]...)
		q.head = 0
	}

	return item
}

// Peek returns the front item without removing it.
func (q *Queue[T]) Peek() T {
	return q.items[q.head]
}

// Len returns the number of items waiting in the queue.
func (q *Queue[T]) Len() int {
	return len(q.items) - q.head
}

// Empty reports whether the queue has no items.
func (q *Queue[T]) Empty() bool {
	return q.Len() == 0
}
