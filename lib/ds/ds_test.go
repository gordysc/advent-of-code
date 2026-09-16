package ds

import "testing"

// TestStackAndQueue checks ordering of the two simple structures.
func TestStackAndQueue(t *testing.T) {
	var s Stack[int]
	s.Push(1, 2, 3)
	if s.Pop() != 3 || s.Peek() != 2 || s.Len() != 2 {
		t.Fatal("Stack wrong")
	}

	q := NewQueue(1, 2, 3)
	if q.Pop() != 1 || q.Peek() != 2 || q.Len() != 2 {
		t.Fatal("Queue wrong")
	}
}

// TestQueueCompaction pushes and pops enough to trigger the compaction path.
func TestQueueCompaction(t *testing.T) {
	q := NewQueue[int]()
	for i := 0; i < 100; i++ {
		q.Push(i)
	}
	for i := 0; i < 100; i++ {
		if got := q.Pop(); got != i {
			t.Fatalf("Pop = %d, want %d", got, i)
		}
	}
	if !q.Empty() {
		t.Fatal("queue should be empty")
	}
}

// TestPriorityQueue checks that the lowest priority pops first.
func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[string]()
	pq.Push("c", 3)
	pq.Push("a", 1)
	pq.Push("b", 2)

	var order string
	for !pq.Empty() {
		v, _ := pq.Pop()
		order += v
	}

	if order != "abc" {
		t.Fatalf("order = %q", order)
	}
}
