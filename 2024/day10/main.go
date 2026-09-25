// Advent of Code 2024, day 10: Hoof It.
// https://adventofcode.com/2024/day/10
package main

import (
	"aoc/lib/grid"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 10, part1, part2)
}

// part1 adds the scores of all trailheads. The score of a trailhead is the
// number of different height 9 positions that its trails reach.
func part1(in string) any {
	topo := parse(in)
	total := 0

	for p, h := range topo.All() {
		if h == 0 {
			total += set.From(trailEnds(topo, p)).Len()
		}
	}

	return total
}

// part2 adds the ratings of all trailheads. The rating of a trailhead is the
// number of different trails that start at it.
func part2(in string) any {
	topo := parse(in)
	total := 0

	for p, h := range topo.All() {
		if h == 0 {
			total += len(trailEnds(topo, p))
		}
	}

	return total
}

// parse reads the topographic map as heights. A cell that is not a digit
// becomes -1, so no trail goes through it.
func parse(in string) grid.Grid[int] {
	return grid.Map(in, func(r rune) int {
		if r < '0' || r > '9' {
			return -1
		}

		return int(r - '0')
	})
}

// trailEnds returns the end of every trail from p to height 9, one entry per
// trail. Each step of a trail goes up, down, left, or right, and increases the
// height by exactly 1. Two trails can end at the same position, so the result
// can hold the same position more than one time.
//
// Because the height increases at each step, a trail can never come back to a
// cell. So the recursion needs no visited set.
func trailEnds(topo grid.Grid[int], p grid.Point) []grid.Point {
	h := topo.At(p)
	if h == 9 {
		return []grid.Point{p}
	}

	var ends []grid.Point
	for _, n := range topo.Neighbors4(p) {
		if topo.At(n) == h+1 {
			// The "..." spreads the slice into separate arguments for append.
			ends = append(ends, trailEnds(topo, n)...)
		}
	}

	return ends
}
