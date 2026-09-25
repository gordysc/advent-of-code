// Advent of Code 2018, day 6: Chronal Coordinates.
// https://adventofcode.com/2018/day/6
package main

import (
	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 6, part1, part2)
}

// limit is the total distance part 2 must stay under in the puzzle. The worked
// example in the puzzle text uses a limit of 32 instead, which gives a region
// of 16. With the real limit, the example region is much larger.
const limit = 10000

// part1 finds the largest area that is closest to one coordinate and does not
// go on forever. An area that reaches the edge of the box around all the
// coordinates keeps going past it, so it is infinite.
func part1(in string) any {
	coords := parse(in)
	lo, hi := bounds(coords)

	area := make([]int, len(coords))
	infinite := make([]bool, len(coords))

	for y := lo.Y; y <= hi.Y; y++ {
		for x := lo.X; x <= hi.X; x++ {
			owner := closest(grid.P(x, y), coords)
			if owner < 0 {
				continue
			}

			area[owner]++
			if x == lo.X || x == hi.X || y == lo.Y || y == hi.Y {
				infinite[owner] = true
			}
		}
	}

	best := 0
	for i, a := range area {
		if !infinite[i] {
			best = max(best, a)
		}
	}

	return best
}

// part2 counts the locations whose total distance to all coordinates is less
// than the limit. Such a location can be outside the box around the
// coordinates, but not by more than limit/len(coords), so the search widens
// the box by that much.
func part2(in string) any {
	coords := parse(in)
	lo, hi := bounds(coords)
	margin := limit / len(coords)

	size := 0

	for y := lo.Y - margin; y <= hi.Y+margin; y++ {
		for x := lo.X - margin; x <= hi.X+margin; x++ {
			p := grid.P(x, y)

			total := 0
			for _, c := range coords {
				total += p.Manhattan(c)
			}

			if total < limit {
				size++
			}
		}
	}

	return size
}

// parse reads one "x, y" coordinate per line.
func parse(in string) []grid.Point {
	var coords []grid.Point

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		coords = append(coords, grid.P(n[0], n[1]))
	}

	return coords
}

// bounds returns the top-left and bottom-right corners of the smallest box
// that holds every coordinate.
func bounds(coords []grid.Point) (grid.Point, grid.Point) {
	lo, hi := coords[0], coords[0]

	for _, c := range coords {
		lo = grid.P(min(lo.X, c.X), min(lo.Y, c.Y))
		hi = grid.P(max(hi.X, c.X), max(hi.Y, c.Y))
	}

	return lo, hi
}

// closest returns the index of the one coordinate nearest to p, or -1 when two
// or more coordinates tie for nearest.
func closest(p grid.Point, coords []grid.Point) int {
	owner, best := -1, -1

	for i, c := range coords {
		d := p.Manhattan(c)

		switch {
		case best < 0 || d < best:
			owner, best = i, d
		case d == best:
			owner = -1
		}
	}

	return owner
}
