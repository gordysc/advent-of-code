package search

import (
	"testing"

	"aoc/lib/grid"
)

// maze is a tiny map used by the search tests.
const maze = `S.#
.##
...
#.E`

// TestBFS walks the maze with breadth-first search.
func TestBFS(t *testing.T) {
	g := grid.Parse(maze)
	start, _ := g.FindByte('S')

	next := func(p grid.Point) []grid.Point {
		var out []grid.Point
		for _, n := range g.Neighbors4(p) {
			if g.At(n) != '#' {
				out = append(out, n)
			}
		}
		return out
	}

	steps, ok := BFS(start, next, func(p grid.Point) bool { return g.At(p) == 'E' })
	if !ok || steps != 5 {
		t.Fatalf("BFS = %d %v, want 5 true", steps, ok)
	}

	if len(Flood(start, next)) != 8 {
		t.Fatalf("Flood reached %d cells, want 8", len(Flood(start, next)))
	}
}

// TestDijkstra uses a weighted triangle where the direct edge is more expensive.
func TestDijkstra(t *testing.T) {
	edges := map[string][]Edge[string]{
		"a": {{"b", 1}, {"c", 10}},
		"b": {{"c", 1}},
	}

	cost, ok := Dijkstra("a", func(s string) []Edge[string] { return edges[s] }, func(s string) bool { return s == "c" })
	if !ok || cost != 2 {
		t.Fatalf("Dijkstra = %d %v", cost, ok)
	}
}

// TestBinarySearch finds the first value passing a threshold.
func TestBinarySearch(t *testing.T) {
	if got := BinarySearch(0, 100, func(n int) bool { return n*n >= 50 }); got != 8 {
		t.Fatalf("BinarySearch = %d, want 8", got)
	}
	if got := BinarySearch(0, 10, func(n int) bool { return false }); got != 11 {
		t.Fatalf("BinarySearch (never) = %d, want 11", got)
	}
}
