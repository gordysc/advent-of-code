// Advent of Code 2022, day 9: Rope Bridge.
// https://adventofcode.com/2022/day/9
//
// The example.txt file holds the first (small) example. It gives 13 for part
// 1 and 1 for part 2. The larger example in the puzzle text gives 36 for
// part 2.
package main

import (
	"strconv"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 9, part1, part2)
}

// part1 counts the positions that the tail of a rope with 2 knots visits.
func part1(in string) any {
	return tailVisits(in, 2)
}

// part2 counts the positions that the tail of a rope with 10 knots visits.
func part2(in string) any {
	return tailVisits(in, 10)
}

// tailVisits moves the head of a rope with the given number of knots, and
// counts the different positions that the last knot visits.
//
// The head moves one step at a time. After each step, each knot follows the
// knot in front of it.
func tailVisits(in string, knots int) int {
	rope := make([]grid.Point, knots)
	visited := set.Of(rope[knots-1])

	for _, line := range input.Lines(in) {
		dir := grid.DirFromRune[rune(line[0])]
		steps, _ := strconv.Atoi(line[2:])

		for range steps {
			rope[0] = rope[0].Add(dir)

			for i := 1; i < knots; i++ {
				rope[i] = follow(rope[i], rope[i-1])
			}

			visited.Add(rope[knots-1])
		}
	}

	return visited.Len()
}

// follow moves a knot one step toward the knot in front of it, when the two
// knots do not touch.
//
// The knots touch when they are at most one step apart, also on a diagonal.
// When they do not touch, the knot moves one step on each axis where the
// positions are different. This covers the straight and the diagonal moves.
func follow(knot, front grid.Point) grid.Point {
	if knot.Chebyshev(front) <= 1 {
		return knot
	}

	diff := front.Sub(knot)

	return knot.Add(grid.P(mathx.Sign(diff.X), mathx.Sign(diff.Y)))
}
