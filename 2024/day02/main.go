// Advent of Code 2024, day 2: Red-Nosed Reports.
// https://adventofcode.com/2024/day/2
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 2, part1, part2)
}

// part1 counts the reports that are safe.
func part1(in string) any {
	count := 0

	for _, line := range input.Lines(in) {
		if safe(input.Ints(line)) {
			count++
		}
	}

	return count
}

// part2 counts the reports that are safe, or that become safe when we remove
// one level.
//
// A report has only a few levels, so we simply try each removal.
func part2(in string) any {
	count := 0

	for _, line := range input.Lines(in) {
		if tolerable(input.Ints(line)) {
			count++
		}
	}

	return count
}

// safe reports whether the levels all increase or all decrease, with each
// step between 1 and 3.
func safe(levels []int) bool {
	if len(levels) < 2 {
		return true
	}

	// The first step sets the direction that all the other steps must follow.
	increasing := levels[1] > levels[0]

	for i := 1; i < len(levels); i++ {
		step := levels[i] - levels[i-1]
		if !increasing {
			step = -step
		}

		if step < 1 || step > 3 {
			return false
		}
	}

	return true
}

// tolerable reports whether the levels are safe, or are safe without one of
// the levels.
func tolerable(levels []int) bool {
	if safe(levels) {
		return true
	}

	for i := range levels {
		// slices.Delete changes the slice it gets, so it works on a copy.
		// slices.Clone makes that copy.
		rest := slices.Delete(slices.Clone(levels), i, i+1)
		if safe(rest) {
			return true
		}
	}

	return false
}
