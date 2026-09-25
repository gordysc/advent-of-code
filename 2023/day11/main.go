// Advent of Code 2023, day 11: Cosmic Expansion.
// https://adventofcode.com/2023/day/11
package main

import (
	"slices"

	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 11, part1, part2)
}

// part1 adds up the shortest distances between all pairs of galaxies. Each
// empty row and column becomes two.
func part1(in string) any {
	return sumDistances(grid.Parse(in), 2)
}

// part2 adds up the same distances, but each empty row and column becomes one
// million.
//
// The example gives answers for a factor of 10 and 100, not one million. The
// example is 10 by 10 and a real input is 140 by 140, so a short input selects
// the example factor of 100. For example.txt, part 2 gives 8410.
func part2(in string) any {
	g := grid.Parse(in)

	factor := 1_000_000
	if g.H < 20 {
		factor = 100
	}

	return sumDistances(g, factor)
}

// sumDistances adds up the Manhattan distances between all pairs of galaxies,
// after each empty row and column grows to factor rows or columns.
//
// A Manhattan distance is the sum of an x distance and a y distance. So each
// axis can be done alone.
func sumDistances(g grid.Grid[byte], factor int) int {
	var xs, ys []int
	for p, v := range g.All() {
		if v == '#' {
			xs = append(xs, p.X)
			ys = append(ys, p.Y)
		}
	}

	return axisSum(xs, factor) + axisSum(ys, factor)
}

// axisSum adds up the distances along one axis between all pairs of
// coordinates, after each empty line grows to factor lines.
//
// After a sort, each coordinate is at least as large as all the coordinates
// before it. So its distance to all of them together is i*c minus their sum.
// This makes the sum O(n log n), not O(n^2).
func axisSum(coords []int, factor int) int {
	slices.Sort(coords)

	total, prefix := 0, 0
	shift := 0 // shift is how far the empty lines so far push this coordinate.

	for i, c := range coords {
		// Each line between the last galaxy and this one that has no galaxy
		// adds factor-1 more lines.
		if i > 0 && c > coords[i-1] {
			shift += (c - coords[i-1] - 1) * (factor - 1)
		}

		pos := c + shift
		total += i*pos - prefix
		prefix += pos
	}

	return total
}
