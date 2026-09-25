// Advent of Code 2023, day 9: Mirage Maintenance.
// https://adventofcode.com/2023/day/9
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 9, part1, part2)
}

// part1 adds up the next value of each history.
func part1(in string) any {
	total := 0

	for _, line := range input.Lines(in) {
		total += next(input.Ints(line))
	}

	return total
}

// part2 adds up the value before the first value of each history.
//
// The value before the start of a history is the next value of the same
// history in reverse order. So part 2 reverses each history and uses the same
// function as part 1.
func part2(in string) any {
	total := 0

	for _, line := range input.Lines(in) {
		values := input.Ints(line)
		slices.Reverse(values)
		total += next(values)
	}

	return total
}

// next finds the value that comes after the last value of a history.
//
// Each difference row gets one more value at its end. The new value of a row
// is its last value plus the new value of the row below it. So the next value
// of the history is the sum of the last values of all the rows.
func next(values []int) int {
	// row is a copy, because the loop changes it in place.
	row := slices.Clone(values)
	sum := 0

	for len(row) > 0 && slices.ContainsFunc(row, func(v int) bool { return v != 0 }) {
		sum += row[len(row)-1]

		// Put the differences into the start of the same slice, then drop
		// the last item, which is not a difference.
		for i := 0; i+1 < len(row); i++ {
			row[i] = row[i+1] - row[i]
		}

		row = row[:len(row)-1]
	}

	return sum
}
