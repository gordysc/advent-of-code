package grid

import (
	"iter"
	"strings"
)

// Grid is a rectangular 2D array stored as one flat slice. Storing cells in a
// single slice is faster and simpler than a slice of slices, and it makes
// copying the grid a one-liner.
type Grid[T any] struct {
	W, H  int
	Cells []T
}

// New returns a w×h grid filled with the zero value of T.
func New[T any](w, h int) Grid[T] {
	return Grid[T]{W: w, H: h, Cells: make([]T, w*h)}
}

// Parse builds a byte grid from puzzle text where every line is one row.
// All lines must have the same length.
func Parse(s string) Grid[byte] {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")

	g := New[byte](len(lines[0]), len(lines))
	for y, line := range lines {
		if len(line) != g.W {
			panic("grid.Parse: rows have different lengths")
		}
		copy(g.Cells[y*g.W:], line)
	}

	return g
}

// Map builds a grid from puzzle text by converting each character with fn.
// Use it for digit grids: grid.Map(s, func(r rune) int { return int(r - '0') }).
func Map[T any](s string, fn func(rune) T) Grid[T] {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")

	g := New[T](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, r := range line {
			g.Cells[y*g.W+x] = fn(r)
		}
	}

	return g
}

// InBounds reports whether p lies inside the grid.
func (g Grid[T]) InBounds(p Point) bool {
	return p.X >= 0 && p.X < g.W && p.Y >= 0 && p.Y < g.H
}

// At returns the cell at p. It panics when p is out of bounds; use [Grid.Get]
// when that may happen.
func (g Grid[T]) At(p Point) T {
	return g.Cells[p.Y*g.W+p.X]
}

// Get returns the cell at p and true, or the zero value and false when p is
// out of bounds.
func (g Grid[T]) Get(p Point) (T, bool) {
	if !g.InBounds(p) {
		var zero T
		return zero, false
	}

	return g.At(p), true
}

// Set writes v into the cell at p.
func (g Grid[T]) Set(p Point, v T) {
	g.Cells[p.Y*g.W+p.X] = v
}

// Points yields every position in reading order (left to right, top to bottom).
// It is an iterator, so use it as: for p := range g.Points() { ... }.
func (g Grid[T]) Points() iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for y := 0; y < g.H; y++ {
			for x := 0; x < g.W; x++ {
				if !yield(Point{x, y}) {
					return
				}
			}
		}
	}
}

// All yields every position together with its value in reading order.
// Use it as: for p, v := range g.All() { ... }.
func (g Grid[T]) All() iter.Seq2[Point, T] {
	return func(yield func(Point, T) bool) {
		for i, v := range g.Cells {
			if !yield(Point{i % g.W, i / g.W}, v) {
				return
			}
		}
	}
}

// Neighbors4 returns the orthogonal neighbours of p that are inside the grid.
func (g Grid[T]) Neighbors4(p Point) []Point {
	return g.inside(p.Neighbors4())
}

// Neighbors8 returns the surrounding points of p that are inside the grid.
func (g Grid[T]) Neighbors8(p Point) []Point {
	return g.inside(p.Neighbors8())
}

// inside filters a list of points down to those in bounds.
func (g Grid[T]) inside(points []Point) []Point {
	out := points[:0]

	for _, q := range points {
		if g.InBounds(q) {
			out = append(out, q)
		}
	}

	return out
}

// Clone returns a deep copy. Grid holds a slice, so plain assignment would share
// the cells between both values.
func (g Grid[T]) Clone() Grid[T] {
	cells := make([]T, len(g.Cells))
	copy(cells, g.Cells)

	return Grid[T]{W: g.W, H: g.H, Cells: cells}
}

// Find returns the first position whose value satisfies pred, in reading order.
func (g Grid[T]) Find(pred func(T) bool) (Point, bool) {
	for i, v := range g.Cells {
		if pred(v) {
			return Point{i % g.W, i / g.W}, true
		}
	}

	return Point{}, false
}

// FindByte returns the first position holding b. It is a shortcut for the very
// common "find the S in the maze" step.
func (g Grid[T]) FindByte(b byte) (Point, bool) {
	// A type switch on the cell type is needed because T is generic and cannot
	// be compared to a byte directly.
	return g.Find(func(v T) bool {
		c, ok := any(v).(byte)
		return ok && c == b
	})
}

// Count returns how many cells satisfy pred.
func (g Grid[T]) Count(pred func(T) bool) int {
	n := 0

	for _, v := range g.Cells {
		if pred(v) {
			n++
		}
	}

	return n
}

// String renders the grid one row per line. Byte cells print as characters,
// anything else uses fmt's %v.
func (g Grid[T]) String() string {
	var sb strings.Builder

	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			v := any(g.Cells[y*g.W+x])
			if b, ok := v.(byte); ok {
				sb.WriteByte(b)
			} else {
				sb.WriteString(strings.TrimSpace(sprint(v)))
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
