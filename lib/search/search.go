// Package search implements the graph searches that dominate Advent of Code:
// breadth-first search, Dijkstra's algorithm and a binary search over integers.
//
// The searches are generic over the state type S. A state can be a grid.Point,
// a struct of position plus direction, or anything else that is comparable,
// because states are stored in maps to remember what has been visited.
package search

import (
	"aoc/lib/ds"
)

// Edge is one move out of a state with its cost, for [Dijkstra].
type Edge[S comparable] struct {
	To   S
	Cost int
}

// BFS finds the smallest number of steps from start to a state where goal is
// true. next returns the states reachable in one step. It returns false when
// the goal cannot be reached.
func BFS[S comparable](start S, next func(S) []S, goal func(S) bool) (int, bool) {
	dist := map[S]int{start: 0}
	queue := ds.NewQueue(start)

	for !queue.Empty() {
		cur := queue.Pop()

		if goal(cur) {
			return dist[cur], true
		}

		for _, n := range next(cur) {
			// A state seen before was reached in fewer or equal steps already.
			if _, seen := dist[n]; seen {
				continue
			}
			dist[n] = dist[cur] + 1
			queue.Push(n)
		}
	}

	return 0, false
}

// Flood runs a breadth-first search from start and returns the distance to
// every reachable state. Use it for "how many cells can I reach" and
// "distance from the start to everywhere" questions.
func Flood[S comparable](start S, next func(S) []S) map[S]int {
	dist := map[S]int{start: 0}
	queue := ds.NewQueue(start)

	for !queue.Empty() {
		cur := queue.Pop()

		for _, n := range next(cur) {
			if _, seen := dist[n]; seen {
				continue
			}
			dist[n] = dist[cur] + 1
			queue.Push(n)
		}
	}

	return dist
}

// Dijkstra finds the cheapest path cost from start to a state where goal is
// true. next returns the outgoing edges with their non-negative costs. It
// returns false when the goal cannot be reached.
func Dijkstra[S comparable](start S, next func(S) []Edge[S], goal func(S) bool) (int, bool) {
	best := map[S]int{start: 0}
	pq := ds.NewPriorityQueue[S]()
	pq.Push(start, 0)

	for !pq.Empty() {
		cur, cost := pq.Pop()

		// The queue can hold stale entries for a state whose cost was later
		// improved. Skipping them is cheaper than removing them from the heap.
		if cost > best[cur] {
			continue
		}

		if goal(cur) {
			return cost, true
		}

		for _, e := range next(cur) {
			newCost := cost + e.Cost
			if old, seen := best[e.To]; seen && old <= newCost {
				continue
			}
			best[e.To] = newCost
			pq.Push(e.To, newCost)
		}
	}

	return 0, false
}

// DijkstraAll returns the cheapest cost from start to every reachable state.
func DijkstraAll[S comparable](start S, next func(S) []Edge[S]) map[S]int {
	best := map[S]int{start: 0}
	pq := ds.NewPriorityQueue[S]()
	pq.Push(start, 0)

	for !pq.Empty() {
		cur, cost := pq.Pop()
		if cost > best[cur] {
			continue
		}

		for _, e := range next(cur) {
			newCost := cost + e.Cost
			if old, seen := best[e.To]; seen && old <= newCost {
				continue
			}
			best[e.To] = newCost
			pq.Push(e.To, newCost)
		}
	}

	return best
}

// BinarySearch returns the smallest n in [lo, hi] for which pred(n) is true.
// pred must be monotonic: false for every value below some threshold and true
// from the threshold on. It returns hi+1 when pred is false for the whole range.
func BinarySearch(lo, hi int, pred func(int) bool) int {
	for lo <= hi {
		mid := lo + (hi-lo)/2

		if pred(mid) {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	return lo
}
