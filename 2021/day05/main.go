// Advent of Code 2021, day 5: Hydrothermal Venture.
// https://adventofcode.com/2021/day/5
package main

import (
	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 5, part1, part2)
}

// part1 counts the points where two or more horizontal or vertical lines
// of vents overlap. It ignores the diagonal lines.
func part1(in string) any {
	return overlaps(in, false)
}

// part2 counts the overlap points again, and now includes the diagonal
// lines. All diagonal lines are at exactly 45 degrees.
func part2(in string) any {
	return overlaps(in, true)
}

// overlaps walks each line one point at a time and counts the points that
// two or more lines cover. The step for each axis is the sign of the
// difference, so the same walk works for horizontal, vertical and 45-degree
// lines. A flat slice of counts is faster than a map for the 1000 by 1000
// area of the puzzle.
func overlaps(in string, diagonals bool) int {
	nums := input.Ints(in)

	// Find the size of the area so the counts fit in a grid.
	w, h := 0, 0

	for i := 0; i+3 < len(nums); i += 4 {
		w = max(w, nums[i]+1, nums[i+2]+1)
		h = max(h, nums[i+1]+1, nums[i+3]+1)
	}

	counts := grid.New[int](w, h)
	result := 0

	for i := 0; i+3 < len(nums); i += 4 {
		from := grid.P(nums[i], nums[i+1])
		to := grid.P(nums[i+2], nums[i+3])

		step := grid.P(mathx.Sign(to.X-from.X), mathx.Sign(to.Y-from.Y))
		if !diagonals && step.X != 0 && step.Y != 0 {
			continue
		}

		for p := from; ; p = p.Add(step) {
			n := counts.At(p) + 1
			counts.Set(p, n)

			// Count each point once, when the second line reaches it.
			if n == 2 {
				result++
			}

			if p == to {
				break
			}
		}
	}

	return result
}
