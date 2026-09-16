package ds

import "container/heap"

// PriorityQueue pops the item with the lowest priority first. It wraps the
// standard library's container/heap, which needs a type that implements
// heap.Interface; that boilerplate lives in the unexported pqHeap type so
// solutions never see it.
type PriorityQueue[T any] struct {
	h *pqHeap[T]
}

// pqItem pairs a value with its priority.
type pqItem[T any] struct {
	value    T
	priority int
}

// pqHeap implements heap.Interface (Len, Less, Swap, Push, Pop) over a slice.
type pqHeap[T any] []pqItem[T]

func (h pqHeap[T]) Len() int           { return len(h) }
func (h pqHeap[T]) Less(i, j int) bool { return h[i].priority < h[j].priority }
func (h pqHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push appends an item; heap.Push calls this and then restores the heap order.
// The pointer receiver is required because it must change the slice header.
func (h *pqHeap[T]) Push(x any) {
	*h = append(*h, x.(pqItem[T]))
}

// Pop removes the last item; heap.Pop first swaps the smallest item to the end.
func (h *pqHeap[T]) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]

	return item
}

// NewPriorityQueue returns an empty min-priority queue.
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{h: &pqHeap[T]{}}
}

// Push adds a value with the given priority. Lower priorities pop first.
func (pq *PriorityQueue[T]) Push(value T, priority int) {
	heap.Push(pq.h, pqItem[T]{value: value, priority: priority})
}

// Pop removes and returns the value with the lowest priority, together with that
// priority. It panics when the queue is empty; check [PriorityQueue.Len] first.
func (pq *PriorityQueue[T]) Pop() (T, int) {
	item := heap.Pop(pq.h).(pqItem[T])

	return item.value, item.priority
}

// Len returns the number of items in the queue.
func (pq *PriorityQueue[T]) Len() int {
	return pq.h.Len()
}

// Empty reports whether the queue has no items.
func (pq *PriorityQueue[T]) Empty() bool {
	return pq.h.Len() == 0
}
