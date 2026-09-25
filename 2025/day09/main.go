// Advent of Code 2025, day 9: Movie Theater.
// https://adventofcode.com/2025/day/9
package main

import (
	"slices"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 9, part1, part2)
}

// part1 returns the largest area of a rectangle that has two red tiles as
// opposite corners. Any tiles can be inside the rectangle.
func part1(in string) any {
	reds := parse(in)
	if len(reds) < 2 {
		return nil
	}

	best := 0

	for i, a := range reds {
		for _, b := range reds[i+1:] {
			best = max(best, area(a, b))
		}
	}

	return best
}

// part2 returns the largest area of a rectangle that has two red tiles as
// opposite corners, and that has only red or green tiles inside it.
//
// The red tiles are the corners of a polygon with only horizontal and vertical
// sides. The green tiles are the sides and the inside of that polygon. Real
// inputs have coordinates near 100000, so a grid of real tiles is too large.
// Instead we compress the coordinates. We keep one column for each x value of a
// red tile, and one column for each gap between two such x values. We do the
// same for the rows. All the tiles in one compressed cell are either all inside
// the polygon or all outside, because no side of the polygon passes through a
// gap. Then a 2D prefix sum tells in constant time if a rectangle holds an
// outside cell.
func part2(in string) any {
	reds := parse(in)
	if len(reds) < 2 {
		return nil
	}

	xs := distinct(reds, func(p grid.Point) int { return p.X })
	ys := distinct(reds, func(p grid.Point) int { return p.Y })

	// Value xs[i] goes to column 2i+1. Column 2i+2 is the gap after it.
	// Column 0 and the last column are a border of empty cells, so that the
	// flood fill below can go around the whole polygon.
	compress := func(p grid.Point) grid.Point {
		cx, _ := slices.BinarySearch(xs, p.X)
		cy, _ := slices.BinarySearch(ys, p.Y)

		return grid.P(2*cx+1, 2*cy+1)
	}

	cells := make([]grid.Point, len(reds))
	for i, p := range reds {
		cells[i] = compress(p)
	}

	outside := outsideCells(cells, 2*len(xs)+1, 2*len(ys)+1)
	sum := prefixSum(outside)

	best := 0

	for i, a := range reds {
		for j := i + 1; j < len(reds); j++ {
			b := reds[j]

			// Test the cheap area first, and skip the pair if it cannot win.
			size := area(a, b)
			if size <= best {
				continue
			}

			if sum.count(cells[i], cells[j]) == 0 {
				best = size
			}
		}
	}

	return best
}

// parse reads one red tile per line, as "x,y".
func parse(in string) []grid.Point {
	var reds []grid.Point

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		if len(n) == 2 {
			reds = append(reds, grid.P(n[0], n[1]))
		}
	}

	return reds
}

// area returns the number of tiles in the rectangle with corners a and b.
// Both corner tiles are part of the rectangle, so each side is one longer
// than the difference of the coordinates.
func area(a, b grid.Point) int {
	return (mathx.Abs(a.X-b.X) + 1) * (mathx.Abs(a.Y-b.Y) + 1)
}

// distinct returns the sorted values that key gives for the points, without
// duplicates.
func distinct(points []grid.Point, key func(grid.Point) int) []int {
	values := make([]int, len(points))
	for i, p := range points {
		values[i] = key(p)
	}

	slices.Sort(values)

	// slices.Compact removes runs of equal values. It needs sorted input to
	// remove all duplicates.
	return slices.Compact(values)
}

// outsideCells marks the compressed cells that are outside the polygon. The
// corners are given in compressed coordinates, in polygon order.
//
// First it draws the sides of the polygon. Then it floods from the corner of
// the grid, which is always outside, and stops at the sides. Every cell that
// the flood does not reach is a side or inside the polygon.
func outsideCells(corners []grid.Point, w, h int) grid.Grid[bool] {
	wall := grid.New[bool](w, h)

	for i, a := range corners {
		// The last corner joins back to the first one.
		b := corners[(i+1)%len(corners)]
		step := grid.P(mathx.Sign(b.X-a.X), mathx.Sign(b.Y-a.Y))

		for p := a; p != b; p = p.Add(step) {
			wall.Set(p, true)
		}

		wall.Set(b, true)
	}

	outside := grid.New[bool](w, h)
	outside.Set(grid.P(0, 0), true)
	queue := []grid.Point{grid.P(0, 0)}

	for len(queue) > 0 {
		p := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		for _, q := range outside.Neighbors4(p) {
			if wall.At(q) || outside.At(q) {
				continue
			}

			outside.Set(q, true)
			queue = append(queue, q)
		}
	}

	return outside
}

// prefix is a 2D prefix sum. Cell (x, y) holds the number of marked cells in
// the rectangle from (0, 0) up to but not including (x, y). So it has one
// more row and column than the grid it counts.
type prefix grid.Grid[int]

// prefixSum builds the 2D prefix sum of the marked cells in g.
func prefixSum(g grid.Grid[bool]) prefix {
	s := grid.New[int](g.W+1, g.H+1)

	for y := range g.H {
		for x := range g.W {
			n := 0
			if g.At(grid.P(x, y)) {
				n = 1
			}

			// Add the cell to the sums above and to the left. The sum above
			// and to the left is in both, so take it away once.
			above := s.At(grid.P(x+1, y))
			left := s.At(grid.P(x, y+1))
			corner := s.At(grid.P(x, y))
			s.Set(grid.P(x+1, y+1), n+above+left-corner)
		}
	}

	return prefix(s)
}

// count returns the number of marked cells in the rectangle with corners a
// and b. Both corners are part of the rectangle.
func (s prefix) count(a, b grid.Point) int {
	g := grid.Grid[int](s)
	x0, x1 := min(a.X, b.X), max(a.X, b.X)+1
	y0, y1 := min(a.Y, b.Y), max(a.Y, b.Y)+1

	return g.At(grid.P(x1, y1)) - g.At(grid.P(x0, y1)) - g.At(grid.P(x1, y0)) + g.At(grid.P(x0, y0))
}
