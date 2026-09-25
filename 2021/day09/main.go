// Advent of Code 2021, day 9: Smoke Basin.
// https://adventofcode.com/2021/day/9
package main

import (
	"slices"

	"aoc/lib/grid"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 9, part1, part2)
}

// part1 adds the risk levels of all low points.
//
// A low point is lower than all of its orthogonal neighbours. Its risk level
// is its height plus 1.
func part1(in string) any {
	heights := parse(in)

	risk := 0
	for _, p := range lowPoints(heights) {
		risk += heights.At(p) + 1
	}

	return risk
}

// part2 multiplies the sizes of the three largest basins.
//
// Each basin flows down to one low point, and height 9 is never part of a
// basin. So a flood fill from each low point, which stops at 9, finds the
// full basin.
func part2(in string) any {
	heights := parse(in)

	next := func(p grid.Point) []grid.Point {
		var out []grid.Point

		for _, q := range heights.Neighbors4(p) {
			if heights.At(q) != 9 {
				out = append(out, q)
			}
		}

		return out
	}

	var sizes []int
	for _, low := range lowPoints(heights) {
		sizes = append(sizes, len(search.Flood(low, next)))
	}

	slices.Sort(sizes)
	slices.Reverse(sizes)

	return sizes[0] * sizes[1] * sizes[2]
}

// lowPoints returns the points that are lower than all their neighbours.
func lowPoints(heights grid.Grid[int]) []grid.Point {
	var lows []grid.Point

	for p, h := range heights.All() {
		isLow := true

		for _, q := range heights.Neighbors4(p) {
			if heights.At(q) <= h {
				isLow = false
				break
			}
		}

		if isLow {
			lows = append(lows, p)
		}
	}

	return lows
}

// parse reads the height map into a grid of digits.
func parse(in string) grid.Grid[int] {
	return grid.Map(in, func(r rune) int { return int(r - '0') })
}
